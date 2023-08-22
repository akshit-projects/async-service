import React from "react";
import {
  Container,
  Typography,
  Grid,
  TextField,
  Select,
  MenuItem,
  Button,
} from "@mui/material";
import { InputLabel, makeStyles } from "@material-ui/core";
const TYPE = "api";
const useStyles = makeStyles((theme) => ({
  formContainer: {
    marginTop: theme.spacing(1),
  },
  selectField: {
    width: "100%",
  },
  inputField: {
    width: "100%",
    marginTop: "8px",
  },
  submitButton: {
    marginTop: "2em",
  },
}));

const APIStep = (props) => {
  const classes = useStyles();
  const setStepProps = props.onUpdateStepProps;
  const stepProps = props.stepProps;
  const onValueChange = (e, name) => {
    const copy = { ...stepProps };
    copy[TYPE][name] = e.target.value;
    setStepProps(copy);
  };

  return (
    <Container sx={{ marginTop: "1em", padding: "0 !important" }}>
      <form>
        <Grid item xs={12} sm={6}>
          <TextField
            label="URL"
            variant="outlined"
            value={stepProps[TYPE]?.url}
            onChange={(e) => onValueChange(e, "url")}
            className={classes.inputField}
          />
        </Grid>
        <Grid container spacing={3} className={classes.formContainer}
            sx={{marginTop: '8px'}}>
          <Grid item xs={12} sm={6}>
            <InputLabel>API Method</InputLabel>
            <Select
              label="Select Option"
              className={classes.selectField}
              variant="outlined"
              value={stepProps[TYPE].method}
              onChange={(e) => onValueChange(e, "method")}
            >
              <MenuItem value="GET">GET</MenuItem>
              <MenuItem value="POST">POST</MenuItem>
              <MenuItem value="PUT">PUT</MenuItem>
            </Select>
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Expected Status"
              variant="outlined"
              value={stepProps[TYPE].expectedStatus}
              sx={{ marginTop: "4px", marginBotton: '0 !important' }}
              onChange={(e) => onValueChange(e, "expectedStatus")}
              className={classes.inputField}
            />
          </Grid>
          {stepProps[TYPE].method !== "GET" && (
            <Grid item xs={12}>
              <TextField
                label="Body"
                multiline
                rows={4}
                value={stepProps[TYPE].body}
                variant="outlined"
                onChange={(e) => onValueChange(e, "body")}
                className={classes.inputField}
              />
            </Grid>
          )}
          <Grid item xs={12}>
            <TextField
              label="Expected Response"
              multiline
              rows={4}
              value={stepProps[TYPE].expectedResponse}
              variant="outlined"
              onChange={(e) => onValueChange(e, "expectedResponse")}
              className={classes.inputField}
            />
          </Grid>
        </Grid>
        {/* <Button variant="contained" color="primary" sx={{ marginTop: "8px" }}>
          Submit
        </Button> */}
      </form>
    </Container>
  );
};

export default APIStep;
