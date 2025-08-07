#!/usr/bin/env python3
import subprocess
import re
import os

CRED_FILE = "/etc/wifi-credentials.txt"
WPA_CONF = "/etc/wpa_supplicant/wpa_supplicant.conf"
INTERFACE = "wlan0"


def read_credentials():
    creds = {}
    with open(CRED_FILE, "r") as f:
        for line in f:
            if "=" in line:
                k, v = line.strip().split("=", 1)
                creds[k] = v
    return creds.get("ssid"), creds.get("password")


def network_exists(ssid):
    if not os.path.exists(WPA_CONF):
        return False
    with open(WPA_CONF, "r") as f:
        conf = f.read()
    # Correct curly braces for .format usage
    return (
        re.search(r'network=\{{[^\}}]*ssid="{}"'.format(re.escape(ssid)), conf)
        is not None
    )


def add_network(ssid, password):
    # Use wpa_passphrase to generate secure config
    result = subprocess.run(
        ["wpa_passphrase", ssid, password], capture_output=True, text=True
    )
    config = result.stdout
    # Remove the commented plaintext password line for security
    config = re.sub(r"^\s*#.*\n", "", config, flags=re.MULTILINE)
    # Append network block to wpa_supplicant.conf
    with open(WPA_CONF, "a") as f:
        f.write("\n" + config)
    print(f"Added network {ssid} to {WPA_CONF}")


def connect_wifi(ssid, password):
    if not network_exists(ssid):
        add_network(ssid, password)
    else:
        print(f"Network {ssid} already exists in {WPA_CONF}")
    # Reload wpa_supplicant config
    subprocess.run(["sudo", "wpa_supplicant", "-B", "-i", INTERFACE, "-c", WPA_CONF, "reconfigure"])


if __name__ == "__main__":
    ssid, password = read_credentials()
    if ssid and password:
        connect_wifi(ssid, password)
    else:
        print("SSID or password missing in credentials file.")
