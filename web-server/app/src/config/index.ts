const backendUri = import.meta.env.VITE_BACKEND_URI
const authorizeUri = import.meta.env.VITE_AUTHORIZE_URI

interface Config {
    authorizeUri: string,
    backendUri: string;
    uriPrefix: string;
    userUri: string;
    uploadCsvUri: string;
}

const config: Config = {
    authorizeUri,
    backendUri,
    uriPrefix: import.meta.env.VITE_URI_PREFIX,
    userUri: `${backendUri}/user`,
    uploadCsvUri: `${backendUri}/upload`,
}

export default config
