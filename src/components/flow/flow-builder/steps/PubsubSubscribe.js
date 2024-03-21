import React from "react";
import {
  Container,
  TextField,
} from "@mui/material";
const TYPE = "messages-subscription";
export default function PubsubSubscribe(props) {
  const setStepProps = props.onUpdateStepProps;
  const stepProps = props.stepProps;
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
            label="Subscription Name"
            onChange={(e) => onValueChange(e, "subscriptionName")}
            required
            value={stepProps[TYPE].subscriptionName}
          />
        </form>
      </Container>
    </>
  );
}
