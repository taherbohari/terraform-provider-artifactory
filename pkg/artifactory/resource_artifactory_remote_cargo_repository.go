package artifactory

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var cargoRemoteSchema = mergeSchema(baseRemoteSchema, map[string]*schema.Schema{
	"git_registry_url": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Description:  `This is the index url, expected to be a git repository. for remote artifactory use "arturl/git/repokey.git"`,
	},
	"anonymous_access": {
		Type:     schema.TypeBool,
		Optional: true,
		Description: "(On the UI: Anonymous download and search) Cargo client does not send credentials when performing download and search for crates. " +
			"Enable this to allow anonymous access to these resources (only), note that this will override the security anonymous access option.",
	},
	"list_remote_folder_items": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: `(Optional) Lists the items of remote folders in simple and list browsing. The remote content is cached according to the value of the 'Retrieval Cache Period'. Default value is 'false'.`,
	},
})

type CargoRemoteRepo struct {
	RemoteRepositoryBaseParams
	RegistryUrl           string `hcl:"git_registry_url" json:"gitRegistryUrl"`
	AnonymousAccess       bool   `hcl:"anonymous_access" json:"cargoAnonymousAccess"`
	ListRemoteFolderItems bool   `json:"listRemoteFolderItems"`
}

func resourceArtifactoryRemoteCargoRepository() *schema.Resource {
	return mkResourceSchema(cargoRemoteSchema, defaultPacker, unpackCargoRemoteRepo, func() interface{} {
		return &CargoRemoteRepo{
			RemoteRepositoryBaseParams: RemoteRepositoryBaseParams{
				Rclass:      "remote",
				PackageType: "cargo",
			},
		}
	})
}

func unpackCargoRemoteRepo(s *schema.ResourceData) (interface{}, string, error) {
	d := &ResourceData{s}
	repo := CargoRemoteRepo{
		RemoteRepositoryBaseParams: unpackBaseRemoteRepo(s, "cargo"),
		RegistryUrl:                d.getString("git_registry_url", false),
		AnonymousAccess:            d.getBool("anonymous_access", false),
		ListRemoteFolderItems:      d.getBool("list_remote_folder_items", false),
	}
	return repo, repo.Id(), nil
}
