package qlooctl

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/gloo/pkg/bootstrap"
	glooflags "github.com/solo-io/gloo/pkg/bootstrap/flags"
	"github.com/solo-io/gloo/pkg/protoutil"
	qloostorage "github.com/solo-io/qloo/pkg/bootstrap"
	"github.com/solo-io/qloo/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/solo-io/glooctl/pkg/config"
	"strings"
	"github.com/pkg/errors"
	"github.com/solo-io/qloo/pkg/api/types/v1"
)

var Opts bootstrap.Options

var outputFormat string

var RootCmd = &cobra.Command{
	Use:   "qlooctl",
	Short: "Interact with QLoo's storage API from the command line",
	Long: "As QLoo features a storage-based API, direct communication with " +
		"the QLoo server is not necessary. qlooctl simplifies the administration of " +
		"QLoo by providing an easy way to create, read, update, and delete QLoo storage objects.\n\n" +
		"" +
		"The primary concerns of qlooctl are Schemas and ResolverMaps. Schemas contain your GraphQL schema;" +
		" ResolverMaps define how your schema fields are resolved.\n\n" +
		"" +
		"Start by creating a schema using qlooctl schema create --from-file <path/to/your/graphql/schema>",
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format for results")
	glooflags.AddConfigStorageOptionFlags(RootCmd, &Opts)
	glooflags.AddFileFlags(RootCmd, &Opts)
	glooflags.AddKubernetesFlags(RootCmd, &Opts)
	glooflags.AddConsulFlags(RootCmd, &Opts)
	config.LoadConfig(&Opts)
}

func printAsYaml(msg proto.Message) error {
	jsn, err := protoutil.Marshal(msg)
	if err != nil {
		return err
	}
	yam, err := yaml.JSONToYAML(jsn)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", yam)
	return nil
}

func printAsJSON(msg proto.Message) error {
	jsn, err := protoutil.Marshal(msg)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", jsn)
	return nil
}

func printTable(msg proto.Message) error {
	switch obj := msg.(type) {
	case *v1.Schema:
		printSchema(obj)
	case *v1.ResolverMap:
		printResolverMap(obj)
	default:
		return errors.Errorf("unknown type %v", msg)
	}
	return nil
}

func printSchema(schema *v1.Schema) {
	fmt.Printf("%v", schema.Name)
}

func printResolverMap(resolverMap *v1.ResolverMap) {
	fmt.Printf("%v", resolverMap.Name)
}

func Print(msg proto.Message) error {
	switch strings.ToLower(outputFormat) {
	case "yaml":
		return printAsYaml(msg)
	case "json":
		return printAsJSON(msg)
	default:
		return printTable(msg)
	}
}

func MakeClient() (storage.Interface, error) {
	return qloostorage.Bootstrap(Opts)
}
