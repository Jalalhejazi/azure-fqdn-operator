package stub

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/api/core/v1"
)

func readyService(s *v1.Service) {

	// Check for PIP and create FQDN if PIP.
	if len(s.Status.LoadBalancer.Ingress) < 1 {
		fmt.Println("Service Created with Azure FQDN annotaton, pending IP address.")
	} else {
		fmt.Println("Creatig Azure FQDN on PIP for IP: " + s.Status.LoadBalancer.Ingress[0].IP)
		createFQDN(s)
		tagService(s)
	}
}

func createFQDN(s *v1.Service) {

	ip := s.Status.LoadBalancer.Ingress[0].IP
	fqdn := s.Annotations["azure-fqdn"]
	rg := s.Annotations["azure-fqdn-rg"]
	location := s.Annotations["zure-fqdn-location"]

	ctx := context.Background()
	ipClient := getIPClient()

	var ipObject, err = getIPObject(ip, ipClient, s)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	err = future.WaitForCompletion(ctx, ipClient.Client)
	if err != nil {
		log.Fatal(err)
	}
}

func getIPClient() network.PublicIPAddressesClient {
	a, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	ipClient := network.NewPublicIPAddressesClient(os.Getenv("AZURE_SUBSCRIPTION_ID"))
	ipClient.Authorizer = a

	return ipClient
}

func getIPObject(ipName string, ipClient network.PublicIPAddressesClient, s *v1.Service) (ip network.PublicIPAddress, err error) {

	var ipObject network.PublicIPAddress
	ctx := context.Background()

	// TODO - how can I get the RG?
	allIP, err := ipClient.List(ctx, s.Annotations["azure-fqdn-rg"])
	if err != nil {
		log.Fatal(err)
	}

	for _, element := range allIP.Values() {
		if *element.IPAddress == ipName {
			fmt.Println("This is the correct IP object")
			ipObject := element
			return ipObject, nil
		}
	}
	return ipObject, err
}

func tagService(s *v1.Service) {
	s.ObjectMeta.Annotations["azure-fqdn-kill"] = "true"
	err := sdk.Update(s)
	if err != nil {
		log.Fatal(err)
	}
}
