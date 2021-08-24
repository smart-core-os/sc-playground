export interface ServerConfig {
  grpcAddress: string;
  httpAddress: string;
  httpsAddress: string;
  insecure?: boolean;
  selfSigned?: boolean;
  mutualTls?: boolean;
}
