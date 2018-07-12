package stub

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
)

func readyService(o *v1.Service) {

	// Check for PIP and create FQDN if PIP.
	if len(o.Status.LoadBalancer.Ingress) < 1 {
		// Pending IP address, do nothing (oppourtunity)
	} else {
		logrus.Info("Creating FQDN: service/" + o.Name + " PIP/" + o.Status.LoadBalancer.Ingress[0].IP)
		createFQDN(o)
		tagService(o)
	}
}

func createFQDN(o *v1.Service) {

	ip := o.Status.LoadBalancer.Ingress[0].IP
	fqdn := o.Annotations["azure-fqdn-value"]
	rg := o.Annotations["azure-fqdn-rg"]
	location := o.Annotations["azure-fqdn-location"]

	ctx := context.Background()
	ipClient := getIPClient()

	var ipObject, err = getIPObject(ip, ipClient, o)
	if err != nil {
		logrus.Fatal(err)
	}

	future, err := ipClient.CreateOrUpdate(
		ctx,
		rg,
		*ipObject.Name,
		network.PublicIPAddress{
			Name:     to.StringPtr(*ipObject.Name),
			Location: to.StringPtr(location),
			PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   network.IPv4,
				PublicIPAllocationMethod: network.Static,
				DNSSettings: &network.PublicIPAddressDNSSettings{
					DomainNameLabel: &fqdn,
				},
			},
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	err = future.WaitForCompletion(ctx, ipClient.Client)
	if err != nil {
		logrus.Fatal(err)
	}
}

func getIPClient() network.PublicIPAddressesClient {
	a, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		logrus.Fatal(err)
	}

	ipClient := network.NewPublicIPAddressesClient(os.Getenv("AZURE_SUBSCRIPTION_ID"))
	ipClient.Authorizer = a

	return ipClient
}

func getIPObject(ipName string, ipClient network.PublicIPAddressesClient, o *v1.Service) (ip network.PublicIPAddress, err error) {

	var ipObject network.PublicIPAddress
	ctx := context.Background()

	// TODO - how can I get the RG?
	allIP, err := ipClient.List(ctx, o.Annotations["azure-fqdn-rg"])
	if err != nil {
		logrus.Fatal(err)
	}

	for _, element := range allIP.Values() {
		if *element.IPAddress == ipName {
			ipObject := element
			return ipObject, nil
		}
	}
	return ipObject, err
}

func tagService(o *v1.Service) {
	o.ObjectMeta.Annotations["azure-fqdn-kill"] = "true"
	err := sdk.Update(o)
	if err != nil {
		logrus.Fatal(err)
	}
}
