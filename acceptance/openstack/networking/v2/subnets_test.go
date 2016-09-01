// +build acceptance networking

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func TestSubnetsList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := subnets.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list subnets: %v", err)
	}

	allSubnets, err := subnets.ExtractSubnets(allPages)
	if err != nil {
		t.Fatalf("Unable to extract subnets: %v", err)
	}

	for _, subnet := range allSubnets {
		PrintSubnet(t, &subnet)
	}
}

func TestSubnetCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	// Create Network
	network, err := CreateNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create network: %v", err)
	}
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	if err != nil {
		t.Fatalf("Unable to create subnet: %v", err)
	}
	defer DeleteSubnet(t, client, subnet.ID)

	PrintSubnet(t, subnet)

	// Update Subnet
	newSubnetName := tools.RandomString("TESTACC-", 8)
	updateOpts := subnets.UpdateOpts{
		Name: newSubnetName,
	}
	_, err = subnets.Update(client, subnet.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update subnet: %v", err)
	}

	// Get subnet
	newSubnet, err := subnets.Get(client, subnet.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get subnet: %v", err)
	}

	PrintSubnet(t, newSubnet)
}
