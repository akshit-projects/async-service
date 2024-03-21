import React from "react";
import {
  Container,
  TextField,
} from "@mui/material";
const TYPE = "subscribe-kafka-topic";
export default function KafkaSubscribe(props) {
  const setStepProps = props.onUpdateStepProps;
  const stepProps = props.stepProps;
  console.log('hh', stepProps);
  const onValueChange = (e, name) => {
    const copy = { ...stepProps };
    copy[TYPE][name] = e.target.value;
    setStepProps(copy);
  };

  const onIntValueChange = (e, name) => {
    const copy = { ...stepProps };
    copy[TYPE][name] = parseInt(e.target.value);
    setStepProps(copy);
  };

  const onKafkaConfigValueChange = (e) => {
    console.log('chjan')
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
            variant="outlined"
            label="Max Messages"
            onChange={(e) => onIntValueChange(e, "maxMessages")}
            required
            value={stepProps[TYPE].maxMessages}
          />
          <TextField
            fullWidth
            sx={{ marginTop: "8px" }}
            variant="outlined"
            label="Group Id"
            onChange={(e) => onValueChange(e, "groupId")}
            required
            value={stepProps[TYPE].groupId}
          />
        </form>
      </Container>
    </>
  );
}
