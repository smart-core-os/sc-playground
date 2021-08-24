let _serverConfig = null;

export async function serverConfig() {
  if (_serverConfig === null) {
    /** @type {string} */
    const url = import.meta.env.VITE_CONFIG_URL || '/__/playground/config.json';
    console.debug('server config path:', url);
    _serverConfig = fetch(url)
        .then(res => res.json());
  }
  // todo: retry on network failure
  return _serverConfig;
}

export async function grpcWebEndpoint() {
  /** @type {ServerConfig} */
  const config = await serverConfig();
  let address = config.insecure ? config.httpAddress : (config.httpsAddress || config.httpAddress)
  if (address[0] === ':') {
    // no host
    address = location.hostname + address;
  }

  const protocol = (config.insecure || !config.httpsAddress) ? 'http://' : 'https://';
  return protocol + address;
}
