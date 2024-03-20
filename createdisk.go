package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/common"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func main() {
	// Create an authorizer from environment variables or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatalf("Failed to create Azure authorizer: %v", err)
	}

	// Create a new DisksClient
	disksClient := v4.NewDisksClient("<subscriptionID>")
	disksClient.Authorizer = authorizer

	// Set parameters for the disk creation
	resourceGroupName := "<resourceGroupName>"
	diskName := "<diskName>"
	snapshotID := "<snapshotID>"
	location := "<location>"

	// Create the data disk from snapshot
	dataDisk := v4.Disk{
		Location: to.StringPtr(location),
		DiskProperties: &v4.DiskProperties{
			CreationData: &v4.CreationData{
				CreateOption: v4.DiskCreateOptionTypesCopy,
				SourceResourceID: to.StringPtr(snapshotID),
			},
		},
	}

	// Create the disk
	_, err := disksClient.CreateOrUpdate(context.Background(), resourceGroupName, diskName, dataDisk)
	if err != nil {
		log.Fatalf("Failed to create disk: %v", err)
	}

	fmt.Printf("Data Disk created: %s\n", diskName)

	// Now, you need to get the VMSS and update its configuration to include the newly created disk.
	// You can use the VMSSClient to achieve this.

	// Create a VMSS client
	vmssClient := v4.NewVirtualMachineScaleSetsClient("<subscriptionID>")
	vmssClient.Authorizer = authorizer

	// Get the VMSS
	vmss, err := vmssClient.Get(context.Background(), resourceGroupName, "<vmssName>")
	if err != nil {
		log.Fatalf("Failed to get VMSS: %v", err)
	}

	// Add the new disk to the VMSS configuration
	vmss.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.DataDisks = append(
		*vmss.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.DataDisks,
		v4.DataDisk{
			Lun:          to.Int32Ptr(1), // Adjust the LUN value as needed
			CreateOption: v4.DiskCreateOptionTypesAttach,
			DiskSizeGB:   to.Int32Ptr(100), // Adjust the disk size as needed
			Name:         to.StringPtr(diskName),
			ManagedDisk: &v4.ManagedDiskParameters{
				ID: to.StringPtr(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/disks/%s", "<subscriptionID>", resourceGroupName, diskName)),
			},
		},
	)

	// Update the VMSS with the modified configuration
	_, err = vmssClient.CreateOrUpdate(context.Background(), resourceGroupName, "<vmssName>", vmss)
	if err != nil {
		log.Fatalf("Failed to update VMSS: %v", err)
	}

	fmt.Println("Data Disk attached to VMSS successfully.")
}
