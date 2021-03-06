package logicmonitor

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
)

func TestAccLogicMonitorCollector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogicMonitorConfigCollector,
				Check: resource.ComposeTestCheckFunc(
					testCollectorExists("logicmonitor_collector.collector1"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "description", "test collector"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "enable_collector_device_failover", "false"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "enable_failback", "false"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "resend_interval", "5"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector1", "suppress_alert_clear", "false"),

					testCollectorExists("logicmonitor_collector.collector2"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector2", "description", "test collector2"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector2", "enable_collector_device_failover", "true"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector2", "enable_failback", "true"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector2", "resend_interval", "1"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector2", "suppress_alert_clear", "true"),
					resource.TestCheckResourceAttr(
						"logicmonitor_collector.collector2", "escalation_chain_id", "1"),
				),
			},
		},
	})
}

func testCollectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*lmv1.DefaultApi)
	if err := testCollectorDestroyHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testCollectorDestroyHelper(s *terraform.State, client *lmv1.DefaultApi) error {
	for _, r := range s.RootModule().Resources {
		id, e := strconv.Atoi(r.Primary.ID)
		if e != nil {
			return e
		}

		restCollectorResponse, apiResponse, e := client.DeleteCollectorById(int32(id))
		err := checkStatus(restCollectorResponse.Status, restCollectorResponse.Errmsg, apiResponse.StatusCode, apiResponse.Status, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func testCollectorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*lmv1.DefaultApi)
		if err := testCollectorExistsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func testCollectorExistsHelper(s *terraform.State, client *lmv1.DefaultApi) error {
	collectorID := s.RootModule().Resources["logicmonitor_collector.collector1"]
	id, e := strconv.Atoi(collectorID.Primary.ID)
	if e != nil {
		return e
	}

	restCollectorResponse, apiResponse, e := client.GetCollectorById(int32(id), "")
	err := checkStatus(restCollectorResponse.Status, restCollectorResponse.Errmsg, apiResponse.StatusCode, apiResponse.Status, e)
	if err != nil {
		return err
	}
	return nil
}

const testAccCheckLogicMonitorConfigCollector = `
resource "logicmonitor_collector" "collector1" {
    description                       = "test collector"
    enable_collector_device_failover  = false
    enable_failback                   = false
    resend_interval                   = 5
    suppress_alert_clear              = false
}

resource "logicmonitor_collector" "collector2" {
    description                       = "test collector2"
    enable_collector_device_failover  = true
    enable_failback                   = true
    resend_interval                   = 1
    suppress_alert_clear              = true
    escalation_chain_id               = 1

    backup_collector_id = 16
}
`
