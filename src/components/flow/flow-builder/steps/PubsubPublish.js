import React from "react";
import { TextField, Container } from "@mui/material";
const TYPE = "publish-message";
export default function PubsubPublish(props) {
  const setStepProps = props.onUpdateStepProps;
  const stepProps = props.stepProps;

  const onMessageChange = (e) => {
    const copy = { ...stepProps };
    copy[TYPE].messages = [e.target.value];
    setStepProps(copy);
  };

  const onValueChange = (e, name) => {
    const copy = { ...stepProps };
    copy[TYPE][name] = e.target.value;
    setStepProps(copy);
  };

  return (
    <>
      <Container sx={{ marginTop: "1em", padding: "0 !important" }}>
        <form>
          <TextField
            fullWidth
            variant="outlined"
            label="Project ID"
            onChange={(e) => onValueChange(e, "projectId")}
            required
            value={stepProps[TYPE].projectId}
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
            value={stepProps[TYPE].messages[0]}
          />
        </form>
      </Container>
    </>
  );
}
