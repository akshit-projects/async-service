module.exports = {
    // API
    PATHS: {
        FLOWS: '/flow',
        ADD_FLOW: '/flow/new',
        OPEN_FLOW: '/flow/:id',
        SUITES: '/suite',
    },

    // FLOW
    DEFAULT_FLOW_NAME: 'Untitled',
    FLOW_FUNCTIONS: {
        API: 'api',
        PUBLISH_MESSAGE: 'publish-message',
        MESSAGES_SUBSCRIPTION: 'messages-subscription',
        PUBLISH_KAFKA_MESSAGE: 'publish-kafka-message',
        SUBSCRIBE_KAFKA_TOPIC: 'subscribe-kafka-topic'
    },
    FLOW_RESPONSE_STATES: {
        PROGRESS: 'PROGRESS',
        ERROR: 'error',
        SUCCESS: 'SUCCESS',
    },
    FLOW_NEW_PATH_SUFFIX: 'new',


    // OTHERS
    BACKEND_URL: process.env.REACT_APP_BACKEND_URL,
    WS_BACKEND_URL: process.env.REACT_APP_WS_BACKEND_URL,
    HTTP_STATUS_RATE_LIMIT_EXCEEDED: 429,
};