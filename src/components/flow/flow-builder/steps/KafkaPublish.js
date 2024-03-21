import React from "react";
import { TextField, Container } from "@mui/material";
const TYPE = "publish-kafka-message";
export default function KafkaPublish(props) {
  const setStepProps = props.onUpdateStepProps;
  const stepProps = props.stepProps;
  console.log(stepProps);
  const onMessageChange = (e) => {
    const copy = { ...stepProps };
    copy[TYPE].messages = [ { value: e.target.value } ];
    setStepProps(copy);
  };

  const onValueChange = (e, name) => {
    const copy = { ...stepProps };
    copy[TYPE][name] = e.target.value;
    setStepProps(copy);
  };

  const onKafkaConfigValueChange = (e) => {
    const copy = { ...stepProps };
    copy[TYPE]['kafkaConfig'] = {
        'bootstrapServers': [ 'localhost:9092' ]
    };
    setStepProps(copy);
  };

  return (
    <>
      <Container sx={{ marginTop: "1em", padding: "0 !important" }}>
        <form>
          <TextField
            fullWidth
            variant="outlined"
            label="Kafka Cluster"
            onChange={onKafkaConfigValueChange}
            required
            value={stepProps[TYPE].kafkaConfig}
          />
          <TextField
            fullWidth
            sx={{ marginTop: "8px" }}
            variant="outlined"
            label="Topic Name"
            onChange={(e) => onValueChange(e, "topicName")}
            required
            value={stepProps[TYPE].topicName}
          />

          <TextField
            fullWidth
            sx={{ marginTop: "8px" }}
            label="Message"
            variant="outlined"
            onChange={onMessageChange}
            multiline
            rows={6}
            required
            value={stepProps[TYPE].messages?.[0]?.value}
          />
        </form>
      </Container>
    </>
  );
}
