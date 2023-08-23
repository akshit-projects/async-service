import { TextField, Button } from "@mui/material";
import React from "react";
import Modal from "react-modal";

Modal.setAppElement("#root");
export default function SuccessModal(props) {
  const isModalOpen = props.isModalOpen;
  const closeRequest = props.closeRequest;
  const value = JSON.stringify(props.value?.response || '');

  return (
    <Modal
      isOpen={isModalOpen}
      className="modal"
      overlayClassName="overlay"
    >
      <Button
        onClick={closeRequest}
        sx={{ padding: "0 !important", minWidth: "unset" }}
      >
        <i className="material-icons">close</i>
      </Button>
      <TextField fullWidth value={value} multiline rows={10} />
    </Modal>
  );
}
