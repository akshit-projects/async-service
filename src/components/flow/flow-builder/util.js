import constants from "../../../constants/constants";

const validMethods = ["GET", "POST", "PUT", "DELETE"];

export const validateSteps = (steps) => {
  if (steps.length === 0) {
    return new Error("Minimum one step is required");
  }
  for (let idx = 0; idx < steps.length; idx++) {
    const step = steps[idx];
    if (!step.name) {
      return new Error("Name is required for step " + (idx + 1));
    }
    const err = validateStep(step);
    if (err) {
      return new Error("Step: " + (idx + 1) + " --> " + err.message);
    }
  }
  return null;
}

function validateStep(step) {
  switch (step.function) {
    case constants.FLOW_FUNCTIONS.API:
      return validateHttpStep(step);
    // case constants.FLOW_FUNCTIONS.PUBLISH_MESSAGE:
    //   return validatePubsubPublish(step);
    // case constants.FLOW_FUNCTIONS.MESSAGES_SUBSCRIPTION:
    //   return validatePubsubSubscribe(step);
    case constants.FLOW_FUNCTIONS.PUBLISH_KAFKA_MESSAGE:
      return validateKafkaPublishMessage(step);
      case constants.FLOW_FUNCTIONS.SUBSCRIBE_KAFKA_TOPIC:
        return validateKafkaSubscribeTopic(step);
    default:
      return new Error("Invalid step function");
  }
}

// Pubsub subscribe validation block
function validatePubsubSubscribe(step) {
  const subscribeRequest = step.meta;

  if (!subscribeRequest) {
    return new Error("Unable to get subscribe request data");
  }

  if (!subscribeRequest.projectId) {
    return new Error("Project id is required for subscription step");
  }

  if (!subscribeRequest.subscriptionName) {
    return new Error("Subscription name is required for subscription step");
  }

  return null;
}

// Pubsub publish validation block
function validatePubsubPublish(step) {
  const publishRequest = step.meta;
  if (!publishRequest) {
    return new Error("Unable to get publish request data");
  }

  if (!publishRequest.projectId) {
    return new Error("Project id required for publish request");
  }

  if (!publishRequest.topicName) {
    return new Error("Topic name is required for publish request");
  }

  if (publishRequest.messages.length === 0) {
    return new Error("At least one message is required for publish request");
  }

  return null;
}

// Kafka publish validation block
function validateKafkaPublishMessage(step) {
    const publishRequest = step.meta;
    console.log(publishRequest);
  
    if (!publishRequest) {
      return new Error("Unable to get kafka publish message request data");
    }
  
    if (!publishRequest.kafkaConfig) {
      return new Error("Kafka cluster required for publishing message.");
    }
  
    if (!publishRequest.topicName) {
      return new Error("Topic name is required for kafka publish request");
    }
  
    if (publishRequest.messages.length === 0) {
      return new Error("At least one message is required for kafka publish request");
    }
  
    return null;
}


// Kafka subscription validation block
function validateKafkaSubscribeTopic(step) {
    const subscribeRequest = step.meta;
    console.log(subscribeRequest);
  
    if (!subscribeRequest) {
      return new Error("Unable to get kafka publish message request data");
    }
  
    if (!subscribeRequest.kafkaConfig) {
      return new Error("Kafka cluster required for publishing message.");
    }
  
    if (!subscribeRequest.topicName) {
      return new Error("Topic name is required for kafka subscribe request");
    }
  
    return null;
  }

// HTTP validation block
function validateHttpStep(step) {
  const httpReq = step.meta;

  if (!httpReq) {
    return new Error("Unable to get http request data");
  }

  httpReq.method = httpReq.method?.toUpperCase();

  const methodValidationError = validateHTTPMethod(httpReq.method);
  if (methodValidationError) {
    return methodValidationError;
  }

  if (httpReq.method === "GET" && httpReq.body) {
    return new Error("Body can't go with GET method");
  } else if (httpReq.method !== "GET" && httpReq.body === null) {
    return new Error("Body is required for " + httpReq.method);
  }

  try {
    new URL(httpReq.url);
  } catch (err) {
    return new Error("Invalid request URL passed");
  }

  return null;
}

function validateHTTPMethod(method) {
  if (!validMethods.includes(method)) {
    return new Error("Invalid http method provided");
  }
  return null;
}
