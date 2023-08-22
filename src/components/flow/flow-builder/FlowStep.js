import React, { useEffect, useState } from "react";
import { FormControl, InputLabel, MenuItem, Select } from "@mui/material";

import APIStep from "./APIStep";
import PubsubPublish from "./PubsubPublish";
import PubsubSubscribe from "./PubsubSubscribe";

const FlowStep = ({ step, index, onUpdate }) => {
    const [stepProps, setStepProps] = useState({
        api: step?.value?.api || {
            method: "GET",
            type: "REST"
        },
        "publish-message": step?.value?.["publish-message"] || {
            projectId: "",
            topicName: "",
            messages: [],
            type: "pubsub",
        },
        "messages-subscription": step?.value?.["messages-subscription"] || {
            projectId: "",
            subscriptionName: "",
            type: "pubsub",
        },
    }
    );
    console.log(index, stepProps, step);
    
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
            <MenuItem value="api">API Request</MenuItem>
            <MenuItem value="publish-message">Publish Message</MenuItem>
            <MenuItem value="messages-subscription">Subscribe Messages</MenuItem>
          </Select>
        </FormControl>
        {step.functionType === "api" && (
          <>
            <APIStep onUpdateStepProps={setStepProps} stepProps={stepProps} />
          </>
        )}
        {step.functionType === "publish-message" && (
          <PubsubPublish
            onUpdateStepProps={setStepProps}
            stepProps={stepProps}
          />
        )}
        {step.functionType === "messages-subscription" && (
          <PubsubSubscribe
            onUpdateStepProps={setStepProps}
            stepProps={stepProps}
          />
        )}
      </div>
    </div>
  );
};

export default FlowStep;
