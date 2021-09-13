package genesyscloud

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mypurecloud/platform-client-sdk-go/v46/platformclientv2"
)

func TestAccResourceFlow(t *testing.T) {
	var (
		flowResource1 = "test_flow1"
		flowResource2 = "test_flow2"
		flowName1     = "Terraform Flow Test-" + uuid.NewString()
		flowName2     = "Terraform Flow Test-" + uuid.NewString()
		flowType1     = "INBOUNDCALL"
		flowType2     = "INBOUNDEMAIL"
		filePath1     = "../examples/resources/genesyscloud_architect_flow/inboundcall_flow_example.yaml"
		filePath2     = "../examples/resources/genesyscloud_architect_flow/inboundcall_flow_example2.yaml"
		filePath3     = "../examples/resources/genesyscloud_architect_flow/inboundcall_flow_example3.yaml"

		inboundcallConfig1 = fmt.Sprintf("inboundCall:\n  name: %s\n  defaultLanguage: en-us\n  startUpRef: ./menus/menu[mainMenu]\n  initialGreeting:\n    tts: Archy says hi!!!\n  menus:\n    - menu:\n        name: Main Menu\n        audio:\n          tts: You are at the Main Menu, press 9 to disconnect.\n        refId: mainMenu\n        choices:\n          - menuDisconnect:\n              name: Disconnect\n              dtmf: digit_9", flowName1)
		inboundcallConfig2 = fmt.Sprintf("inboundCall:\n  name: %s\n  defaultLanguage: en-us\n  startUpRef: ./menus/menu[mainMenu]\n  initialGreeting:\n    tts: Archy says hi!!!!!\n  menus:\n    - menu:\n        name: Main Menu\n        audio:\n          tts: You are at the Main Menu, press 9 to disconnect.\n        refId: mainMenu\n        choices:\n          - menuDisconnect:\n              name: Disconnect\n              dtmf: digit_9", flowName2)

		inboundemailConfig1 = fmt.Sprintf(`inboundEmail:
    name: %s
    division: Home
    startUpRef: "/inboundEmail/states/state[Initial State_10]"
    defaultLanguage: en-us
    supportedLanguages:
        en-us:
            defaultLanguageSkill:
                noValue: true
    settingsInboundEmailHandling:
        emailHandling:
            disconnect:
                none: true
    settingsErrorHandling:
        errorHandling:
            disconnect:
                none: true
    states:
        - state:
            name: Initial State
            refId: Initial State_10
            actions:
                - disconnect:
                    name: Disconnect
`, flowName1)
	)

	os.Setenv("GENESYSCLOUD_OAUTHCLIENT_ID", "df4cf7c9-bdcd-4c87-bb90-969455486dd1")
	os.Setenv("GENESYSCLOUD_OAUTHCLIENT_SECRET", "1zjnIHkin-5UKH_u2dLbHsoax6K9kvj0ZNhi8wHJY6w")
	os.Setenv("GENESYSCLOUD_REGION", "dca")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Create flow
				Config: generateFlowResource(
					flowResource1,
					flowName1,
					flowType1,
					filePath1,
					falseValue,
					trueValue,
					trueValue,
					inboundcallConfig1,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource1, "description", fmt.Sprintf("Flow name: %s, Flow type: %s", flowName1, flowType1)),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource1, "debug", falseValue),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource1, "force_unlock", trueValue),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource1, "recreate", trueValue),
				),
			},
			{
				// Update flow with name
				Config: generateFlowResource(
					flowResource1,
					flowName2,
					flowType1,
					filePath2,
					falseValue,
					trueValue,
					trueValue,
					inboundcallConfig2,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource1, "description", fmt.Sprintf("Flow name: %s, Flow type: %s", flowName2, flowType1)),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource1, "debug", falseValue),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource1, "force_unlock", trueValue),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource1, "recreate", trueValue),
				),
			},
			{
				// Import/Read
				ResourceName:            "genesyscloud_architect_flow." + flowResource1,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"filepath", "debug", "force_unlock", "recreate"},
			},
			{
				// Create inboundemail flow
				Config: generateFlowResource(
					flowResource2,
					flowName1,
					flowType2,
					filePath3,
					falseValue,
					trueValue,
					trueValue,
					inboundemailConfig1,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource2, "description", fmt.Sprintf("Flow name: %s, Flow type: %s", flowName1, flowType2)),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource2, "debug", falseValue),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource2, "force_unlock", trueValue),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource2, "recreate", trueValue),
				),
			},
			{
				// Update inboundemail flow to inboundcall
				Config: generateFlowResource(
					flowResource2,
					flowName2,
					flowType1,
					filePath2,
					falseValue,
					trueValue,
					trueValue,
					inboundcallConfig2,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource2, "description", fmt.Sprintf("Flow name: %s, Flow type: %s", flowName2, flowType1)),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource2, "debug", falseValue),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource2, "force_unlock", trueValue),
					resource.TestCheckResourceAttr("genesyscloud_architect_flow."+flowResource2, "recreate", trueValue),
				),
			},
			{
				// Import/Read
				ResourceName:            "genesyscloud_architect_flow." + flowResource2,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"filepath", "debug", "force_unlock", "recreate"},
			},
		},
		CheckDestroy: testVerifyFlowDestroyed,
	})
}

func generateFlowResource(resourceID string, name string, flowtype string, filepath string, debug string, forceUnlock string, recreate string, filecontent string) string {

	updateFile(filepath, filecontent)

	return fmt.Sprintf(`resource "genesyscloud_architect_flow" "%s" {
        description = "Flow name: %s, Flow type: %s"
        filepath = %s
        debug = %s
        force_unlock = %s
        recreate = %s
	}
	`, resourceID, name, flowtype, strconv.Quote(filepath), debug, forceUnlock, recreate)
}

func updateFile(filepath string, content string) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	file.WriteString(content)

}

func testVerifyFlowDestroyed(state *terraform.State) error {
	architectAPI := platformclientv2.NewArchitectApi()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "genesyscloud_architect_flow" {
			continue
		}

		flow, resp, err := architectAPI.GetFlow(rs.Primary.ID, false)
		if flow != nil {
			return fmt.Errorf("Flow (%s) still exists", rs.Primary.ID)
		} else if resp != nil && resp.StatusCode == 410 {
			// Flow not found as expected
			continue
		} else {
			// Unexpected error
			return fmt.Errorf("Unexpected error: %s", err)
		}
	}
	// Success. All Flows destroyed
	return nil
}