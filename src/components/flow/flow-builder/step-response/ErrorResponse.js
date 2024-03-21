import { Button, TextField, Grid } from "@mui/material";
import Modal from 'react-modal';
import React from "react";

Modal.setAppElement("#root");
export default function ErrorResponse(props) {
  const response = props.response || {};
  const isModalOpen = props.isModalOpen;
  const closeRequest = props.closeRequest;
  console.log(response.response);
  if (!response.actual) {
    const value = JSON.stringify(response.error || '');
    return (
      <Modal isOpen={isModalOpen} className="modal" overlayClassName="overlay">
        <Button
          onClick={closeRequest}
          sx={{ padding: "0 !important", minWidth: "unset" }}
        >
          <i className="material-icons">close</i>
        </Button>
        <TextField fullWidth value={value} sx={{ color: 'red' }} />
      </Modal>
    );
  } else {
    const expected = response.expected;
    const actual = response.actual;
    return (
        <Modal isOpen={isModalOpen} className="modal" overlayClassName="overlay">
          <Button
            onClick={closeRequest}
            sx={{ padding: "0 !important", minWidth: "unset" }}
          >
            <i className="material-icons">close</i>
          </Button>
          <Grid container spacing={2}>
            <Grid item xs={6}>
                <TextField fullWidth value={expected} multiline rows={10} />
            </Grid>
            <Grid item xs={6}>
                <TextField fullWidth value={actual} multiline rows={10} />
            </Grid>
          </Grid>
          <p>{response.error}</p>
        </Modal>
      );
  }
}
