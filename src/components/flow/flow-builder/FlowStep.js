import React, { useEffect, useState } from "react";
import { FormControl, InputLabel, MenuItem, Select } from "@mui/material";

import APIStep from "./steps/APIStep";
import PubsubPublish from "./steps/PubsubPublish";
import PubsubSubscribe from "./steps/PubsubSubscribe";
import constants from "../../../constants/constants";
import KafkaPublish from "./steps/KafkaPublish";
import KafkaSubscribe from "./steps/KafkaSubscribe";

const FlowStep = ({ step, index, onUpdate }) => {
  const [stepProps, setStepProps] = useState({
    [constants.FLOW_FUNCTIONS.API]: step?.value?.api || {
      method: "GET",
      type: "REST",
    },
    [constants.FLOW_FUNCTIONS.PUBLISH_MESSAGE]: step?.value?.[
      constants.FLOW_FUNCTIONS.PUBLISH_MESSAGE
    ] || {
      projectId: "",
      topicName: "",
      messages: [],
      type: "pubsub",
    },
    [constants.FLOW_FUNCTIONS.MESSAGES_SUBSCRIPTION]: step?.value?.[
      constants.FLOW_FUNCTIONS.MESSAGES_SUBSCRIPTION
    ] || {
      projectId: "",
      subscriptionName: "",
      type: "pubsub",
    },
    [constants.FLOW_FUNCTIONS.PUBLISH_KAFKA_MESSAGE]: step?.value?.[
        constants.FLOW_FUNCTIONS.SUBSCRIBE_KAFKA_TOPIC
      ] || {
        kafkaConfig: {},
        topicName: "",
        type: "kafka",
      },
    [constants.FLOW_FUNCTIONS.SUBSCRIBE_KAFKA_TOPIC]: step?.value?.[
        constants.FLOW_FUNCTIONS.PUBLISH_KAFKA_MESSAGE
      ] || {
        kafkaConfig: {},
        topicName: "",
        type: "kafka",
      },
  });

  const handleUpdate = (field, value) => {
    onUpdate(index, field, value);
  };

  useEffect(() => {
    handleUpdate("value", stepProps);
  }, [stepProps]);

  return (
    <div className="step">
      <div>
        <FormControl variant="outlined" fullWidth>
          <InputLabel>Function</InputLabel>
          <Select
            value={step.functionType}
            onChange={(e) => handleUpdate("functionType", e.target.value)}
            label="Function"
          >
            <MenuItem value="">Select Function</MenuItem>
            <MenuItem value={constants.FLOW_FUNCTIONS.API}>API Request</MenuItem>
            {/* <MenuItem value={constants.FLOW_FUNCTIONS.PUBLISH_MESSAGE}>Publish Message</MenuItem>
            <MenuItem value={constants.FLOW_FUNCTIONS.MESSAGES_SUBSCRIPTION}>
              Subscribe Messages
            </MenuItem> */}
            <MenuItem value={constants.FLOW_FUNCTIONS.PUBLISH_KAFKA_MESSAGE}>
              Publish Kafka Messages
            </MenuItem>
            <MenuItem value={constants.FLOW_FUNCTIONS.SUBSCRIBE_KAFKA_TOPIC}>
              Subscribe Kafka Topic
            </MenuItem>
          </Select>
        </FormControl>
        {step.functionType === constants.FLOW_FUNCTIONS.API && (
          <>
            <APIStep onUpdateStepProps={setStepProps} stepProps={stepProps} />
          </>
        )}
        {step.functionType === constants.FLOW_FUNCTIONS.PUBLISH_MESSAGE && (
          <PubsubPublish
            onUpdateStepProps={setStepProps}
            stepProps={stepProps}
          />
        )}
        {step.functionType === constants.FLOW_FUNCTIONS.MESSAGES_SUBSCRIPTION && (
          <PubsubSubscribe
            onUpdateStepProps={setStepProps}
            stepProps={stepProps}
          />
        )}
        {step.functionType === constants.FLOW_FUNCTIONS.PUBLISH_KAFKA_MESSAGE && (
          <KafkaPublish
            onUpdateStepProps={setStepProps}
            stepProps={stepProps}
          />
        )}
        {step.functionType === constants.FLOW_FUNCTIONS.SUBSCRIBE_KAFKA_TOPIC && (
          <KafkaSubscribe
            onUpdateStepProps={setStepProps}
            stepProps={stepProps}
          />
        )}
      </div>
    </div>
  );
};

export default FlowStep;
