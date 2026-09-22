//
//  virtualization_27_arm64.h
//

#pragma once

#import "virtualization_helper.h"
#import <Virtualization/Virtualization.h>

/* macOS 13 API, reused so start options can carry more than one setting */
void setStartUpFromMacOSRecoveryVZMacOSVirtualMachineStartOptions(void *startOptions, bool startUpFromMacOSRecovery);

/* macOS 27 API */
void *newVZMacGuestProvisioningOptions(const char *fullName, const char *username, const char *password, bool logsInAutomatically, bool enablesRemoteLogin);
bool setGuestProvisioningOptionsVZMacOSVirtualMachineStartOptions(void *startOptions, void *provisioningOptions, void **error);
