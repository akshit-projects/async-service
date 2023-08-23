module.exports = {
    // API
    PATHS: {
        FLOWS: '/flow',
        ADD_FLOW: '/flow/new',
        OPEN_FLOW: '/flow/:id',
        API_PREFIX: '/flow',
    },

    // FLOW
    DEFAULT_FLOW_NAME: 'Untitled',
    FLOW_FUNCTIONS: {
        API: 'api',
        PUBLISH_MESSAGE: 'publish-message',
        MESSAGES_SUBSCRIPTION: 'messages-subscription',
    },
    FLOW_RESPONSE_STATES: {
        PROGRESS: 'PROGRESS',
        ERROR: 'ERROR',
        SUCCESS: 'SUCCESS',
    },
    FLOW_NEW_PATH_SUFFIX: 'new',


    // OTHERS
    BACKEND_URL: process.env.REACT_APP_BACKEND_URL,
    WS_BACKEND_URL: process.env.REACT_APP_WS_BACKEND_URL,
    HTTP_STATUS_RATE_LIMIT_EXCEEDED: 429,
};