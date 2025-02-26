package genesyscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mypurecloud/platform-client-sdk-go/v56/platformclientv2"
)

func getAllIdpGeneric(_ context.Context, clientConfig *platformclientv2.Configuration) (ResourceIDMetaMap, diag.Diagnostics) {
	idpAPI := platformclientv2.NewIdentityProviderApiWithConfig(clientConfig)
	resources := make(ResourceIDMetaMap)

	_, resp, getErr := idpAPI.GetIdentityprovidersGeneric()
	if getErr != nil {
		if isStatus404(resp) {
			// Don't export if config doesn't exist
			return resources, nil
		}
		return nil, diag.Errorf("Failed to get IDP Generic: %v", getErr)
	}

	resources["0"] = &ResourceMeta{Name: "generic"}
	return resources, nil
}

func idpGenericExporter() *ResourceExporter {
	return &ResourceExporter{
		GetResourcesFunc: getAllWithPooledClient(getAllIdpGeneric),
		RefAttrs:         map[string]*RefAttrSettings{}, // No references
	}
}

func resourceIdpGeneric() *schema.Resource {
	return &schema.Resource{
		Description: "Genesys Cloud Single Sign-on Generic Identity Provider. See this page for detailed configuration instructions: https://help.mypurecloud.com/articles/add-a-generic-single-sign-on-provider/",

		CreateContext: createWithPooledClient(createIdpGeneric),
		ReadContext:   readWithPooledClient(readIdpGeneric),
		UpdateContext: updateWithPooledClient(updateIdpGeneric),
		DeleteContext: deleteWithPooledClient(deleteIdpGeneric),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the provider.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"certificates": {
				Description: "PEM or DER encoded public X.509 certificates for SAML signature validation.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"issuer_uri": {
				Description: "Issuer URI provided by the provider.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"target_uri": {
				Description: "Target URI provided by the provider.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"relying_party_identifier": {
				Description: "String used to identify Genesys Cloud to the identity provider.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"disabled": {
				Description: "True if Generic provider is disabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"logo_image_data": {
				Description: "Base64 encoded SVG image.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"endpoint_compression": {
				Description: "True if the Genesys Cloud authentication request should be compressed.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"name_identifier_format": {
				Description: "SAML name identifier format. (urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified | urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress | urn:oasis:names:tc:SAML:1.1:nameid-format:X509SubjectName | urn:oasis:names:tc:SAML:1.1:nameid-format:WindowsDomainQualifiedName | urn:oasis:names:tc:SAML:2.0:nameid-format:kerberos | urn:oasis:names:tc:SAML:2.0:nameid-format:entity | urn:oasis:names:tc:SAML:2.0:nameid-format:persistent | urn:oasis:names:tc:SAML:2.0:nameid-format:transient)",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
				ValidateFunc: validation.StringInSlice([]string{
					"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
					"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
					"urn:oasis:names:tc:SAML:1.1:nameid-format:X509SubjectName",
					"urn:oasis:names:tc:SAML:1.1:nameid-format:WindowsDomainQualifiedName",
					"urn:oasis:names:tc:SAML:2.0:nameid-format:kerberos",
					"urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
					"urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
					"urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
				}, false),
			},
		},
	}
}

func createIdpGeneric(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Creating IDP Generic")
	d.SetId("generic")
	return updateIdpGeneric(ctx, d, meta)
}

func readIdpGeneric(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*providerMeta).ClientConfig
	idpAPI := platformclientv2.NewIdentityProviderApiWithConfig(sdkConfig)

	log.Printf("Reading IDP Generic")

	return withRetriesForRead(ctx, 30*time.Second, d, func() *resource.RetryError {
		generic, resp, getErr := idpAPI.GetIdentityprovidersGeneric()
		if getErr != nil {
			if isStatus404(resp) {
				return resource.RetryableError(fmt.Errorf("Failed to read IDP Generic: %s", getErr))
			}
			return resource.NonRetryableError(fmt.Errorf("Failed to read IDP Generic: %s", getErr))
		}

		if generic.Name != nil {
			d.Set("name", *generic.Name)
		} else {
			d.Set("name", nil)
		}

		if generic.Certificate != nil {
			d.Set("certificates", stringListToSet([]string{*generic.Certificate}))
		} else if generic.Certificates != nil {
			d.Set("certificates", stringListToSet(*generic.Certificates))
		} else {
			d.Set("certificates", nil)
		}

		if generic.IssuerURI != nil {
			d.Set("issuer_uri", *generic.IssuerURI)
		} else {
			d.Set("issuer_uri", nil)
		}

		if generic.SsoTargetURI != nil {
			d.Set("target_uri", *generic.SsoTargetURI)
		} else {
			d.Set("target_uri", nil)
		}

		if generic.RelyingPartyIdentifier != nil {
			d.Set("relying_party_identifier", *generic.RelyingPartyIdentifier)
		} else {
			d.Set("relying_party_identifier", nil)
		}

		if generic.Disabled != nil {
			d.Set("disabled", *generic.Disabled)
		} else {
			d.Set("disabled", nil)
		}

		if generic.LogoImageData != nil {
			d.Set("logo_image_data", *generic.LogoImageData)
		} else {
			d.Set("logo_image_data", nil)
		}

		if generic.EndpointCompression != nil {
			d.Set("endpoint_compression", *generic.EndpointCompression)
		} else {
			d.Set("endpoint_compression", nil)
		}

		if generic.NameIdentifierFormat != nil {
			d.Set("name_identifier_format", *generic.NameIdentifierFormat)
		} else {
			d.Set("name_identifier_format", nil)
		}

		log.Printf("Read IDP Generic")
		return nil
	})
}

func updateIdpGeneric(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	issuerUri := d.Get("issuer_uri").(string)
	targetUri := d.Get("target_uri").(string)
	relyingPartyID := d.Get("relying_party_identifier").(string)
	disabled := d.Get("disabled").(bool)
	logoImageData := d.Get("logo_image_data").(string)
	endpointCompression := d.Get("endpoint_compression").(bool)
	nameIdentifierFormat := d.Get("name_identifier_format").(string)

	sdkConfig := meta.(*providerMeta).ClientConfig
	idpAPI := platformclientv2.NewIdentityProviderApiWithConfig(sdkConfig)

	log.Printf("Updating IDP Generic")
	update := platformclientv2.Genericsaml{
		Name:                   &name,
		IssuerURI:              &issuerUri,
		SsoTargetURI:           &targetUri,
		RelyingPartyIdentifier: &relyingPartyID,
		Disabled:               &disabled,
		LogoImageData:          &logoImageData,
		EndpointCompression:    &endpointCompression,
		NameIdentifierFormat:   &nameIdentifierFormat,
	}

	certificates := buildSdkStringList(d, "certificates")
	if certificates != nil {
		if len(*certificates) == 1 {
			update.Certificate = &(*certificates)[0]
		} else {
			update.Certificates = certificates
		}
	}

	_, _, err := idpAPI.PutIdentityprovidersGeneric(update)
	if err != nil {
		return diag.Errorf("Failed to update IDP Generic: %s", err)
	}

	log.Printf("Updated IDP Generic")
	// Give time for public API caches to update
	// It takes a long time with idp resources
	time.Sleep(20 * time.Second)
	return readIdpGeneric(ctx, d, meta)
}

func deleteIdpGeneric(ctx context.Context, _ *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*providerMeta).ClientConfig
	idpAPI := platformclientv2.NewIdentityProviderApiWithConfig(sdkConfig)

	log.Printf("Deleting IDP Generic")
	_, _, err := idpAPI.DeleteIdentityprovidersGeneric()
	if err != nil {
		return diag.Errorf("Failed to delete IDP Generic: %s", err)
	}

	return withRetries(ctx, 30*time.Second, func() *resource.RetryError {
		_, resp, err := idpAPI.GetIdentityprovidersGeneric()
		if err != nil {
			if isStatus404(resp) {
				// IDP Generic deleted
				log.Printf("Deleted IDP Generic")
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting IDP Generic: %s", err))
		}
		return resource.RetryableError(fmt.Errorf("IDP Generic still exists"))
	})
}
