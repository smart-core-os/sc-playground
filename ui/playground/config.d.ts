export interface ServerConfig {
  grpcAddress: string;
  httpAddress: string;
  httpsAddress: string;
  insecure?: boolean;
}
