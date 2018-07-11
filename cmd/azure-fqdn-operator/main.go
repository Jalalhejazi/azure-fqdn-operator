// Some notes:
// This was build with the operator-sdk, however does not use a CRD.
// https://github.com/operator-framework/operator-sdk/issues/326

package main

import (
	"context"
	"runtime"

	stub "github.com/neilpeterson/azure-fqdn-operator/pkg/stub"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()

	sdk.ExposeMetricsPort()

	// Can I remove the resync (0) and instead create a watch on service resources?
	sdk.Watch("v1", "Service", "default", 5)
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}
