//
//  virtualization_27_arm64.m
//

#import "virtualization_27_arm64.h"

void setStartUpFromMacOSRecoveryVZMacOSVirtualMachineStartOptions(void *startOptions, bool startUpFromMacOSRecovery)
{
#ifdef INCLUDE_TARGET_OSX_13
    if (@available(macOS 13, *)) {
        [(VZMacOSVirtualMachineStartOptions *)startOptions setStartUpFromMacOSRecovery:(BOOL)startUpFromMacOSRecovery];
        return;
    }
#endif
    RAISE_UNSUPPORTED_MACOS_EXCEPTION();
}

/*!
 @abstract Provisioning for a macOS guest's first boot after restore: an account, and whether it logs in automatically and accepts SSH.
 @discussion Only a macOS 27 or newer guest honors these; an earlier guest ignores them.
 */
void *newVZMacGuestProvisioningOptions(const char *fullName, const char *username, const char *password, bool logsInAutomatically, bool enablesRemoteLogin)
{
#ifdef INCLUDE_TARGET_OSX_27
    if (@available(macOS 27, *)) {
        VZMacGuestProvisioningOptions *options = [[VZMacGuestProvisioningOptions alloc] init];
        [options setFullName:[NSString stringWithUTF8String:fullName]];
        [options setUsername:[NSString stringWithUTF8String:username]];
        [options setPassword:[NSString stringWithUTF8String:password]];
        [options setLogsInAutomatically:(BOOL)logsInAutomatically];
        [options setEnablesRemoteLogin:(BOOL)enablesRemoteLogin];
        return options;
    }
#endif
    RAISE_UNSUPPORTED_MACOS_EXCEPTION();
}

bool setGuestProvisioningOptionsVZMacOSVirtualMachineStartOptions(void *startOptions, void *provisioningOptions, void **error)
{
#ifdef INCLUDE_TARGET_OSX_27
    if (@available(macOS 27, *)) {
        return (bool)[(VZMacOSVirtualMachineStartOptions *)startOptions
            setGuestProvisioningOptions:(VZMacGuestProvisioningOptions *)provisioningOptions
                                  error:(NSError *_Nullable *_Nullable)error];
    }
#endif
    RAISE_UNSUPPORTED_MACOS_EXCEPTION();
}
