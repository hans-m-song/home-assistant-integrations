export type DiagnoseInternetResponse = {
  WANAccessType: string; // 'Ethernet'
  UpMaxBitRate: string; // '1000'
  ConnectionStatus: string; // 'Connected'
  X_IPv6Address: string; // ''
  LinkStatus: string; // 'Up'
  DuplexMode: string; // 'Full'
  Status: string; // 'Connected'
  X_IPv6Enable: boolean; // false
  X_IPv6DNSServers: string; // ''
  X_IPv6PrefixList: string; // ''
  HasInternetWan: boolean; // true
  ErrReason: string; // 'Success'
  X_IPv6PrefixLength: number; // 0
  X_IPv4Enable: boolean; // true
  DNSServers: string; // ''
  ExternalIPAddress: string; // ''
  X_IPv6AddressingType: string; // 'SLAAC'
  MACAddress: string; // ''
  X_IPv6ConnectionStatus: string; // 'Pending Disconnect'
  Uptime: number; // 65
  DownMaxBitRate: string; // '1000'
  X_IPv6DefaultGateway: string; // ''
  MaxBitRate: string; // '1000'
  StatusCode: string; // 'Connected'
  DefaultGateway: string; // '
};

export type DeviceCountResponse = {
  PrinterNumbers: number; // 0
  UsbNumbers: number; // 0
  UserNumber: number; // 2
  LanActiveNumber: number; // 3
  PhoneNumber: number; // 2
  DatacardNumber: number; // 0
  ActiveDeviceNumbers: number; // 13
};

export type WizardWifiResponse = {
  WifiFrequency: 2 | 5;
  Numbers: string; // "10"
  WifiEnable: boolean; // true;
  WifiSsid: string; // "42-5G"
}[];

export type DeviceinfoResponse = {
  DeviceName: string; // "HG659"
  SerialNumber: string; // "J3N8W17808904958"
  ManufacturerOUI: string; // "00E0FC"
  UpTime: number; // 1034698
  SoftwareVersion: string; // "V100R001C216B112"
  HardwareVersion: string; // "VER.B"
};

export type WandetectResponse = {
  IPv6DefaultGateway: string; // ""
  ConnectionStatus: string; // "Connecting"
  HasInternet: boolean; // true
  BackupStatus: string; // "1"
  DefaultGateway: string; // ""
  PVCResult: string; // ""
  IPv6AddressingType: string; // "SLAAC"
  IPv6PrefixLength: number; // 0
  VDSLLinkStatus: string; // "NoSignal"
  ID: string; // "InternetGatewayDevice.WANDevice.3.WANConnectionDevice.1.WANPPPConnection.1."
  EthernetLinkStatus: string; // "Up"
  ConnectionType: string; // "PPP_Routed"
  IPv6DNSServers: string; // ""
  IPv6PrefixList: string; // ""
  WanResult: string; // "1"
  SearchingStatus: string; // "Finished"
  IPv6ConnectionStatus: string; // "PendingDisconnect"
  DSLLinkStatus: string; // "NoSignal"
  DNSServers: string; // ""
  AccessStatus: string; // "Up"
  IPv6Address: string; // ""
  ExternalIPAddress: string; // ""
  IPv6Enable: boolean; // false
  Uptime: number; // 0
  ErrReason: string; // "ErrConnectFail"
  UMTSLinkStatus: string; // "Down"
  AccessType: string; // "Ethernet"
  Status: string; // "Fault"
  IPv4Enable: boolean; // true
};
