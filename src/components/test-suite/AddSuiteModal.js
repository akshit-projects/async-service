import {
  Alert,
  Button,
  Checkbox,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
} from "@mui/material";
import React, { useState } from "react";
import Modal from "react-modal";
import constants from "../../constants/constants";
import axios from "axios";
import { useNavigate } from "react-router-dom";

export default function AddSuiteModal(props) {
  const isOpen = props.isOpen || false;
  const [suiteName, setSuiteName] = useState("Untitled");
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const [error, setError] = useState(false);
  if (!isOpen) return <></>;

  const addSuite = () => {
    const body = {
      name: suiteName,
      flowIds: flows.map((flow) => flow.id),
    };
    const options = {
      url: `${constants.BACKEND_URL}/api/v1/suite`,
      method: "POST",
      data: body,
    };
    setIsLoading(true);
    axios(options)
      .then(() => {
        navigate("/suite");
      })
      .catch((err) => {
        setError("Unable to save suite");
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  const handleCheckboxChange = props.handleCheckboxChange;
  const closeModal = props.closeModal;

  const flows = props.flows || [];
  return (
    <Modal
      isOpen={isOpen}
      className="modal"
      overlayClassName="overlay"
      sx={{ padding: "24px" }}
    >
      <TextField
        label="Suite name"
        value={suiteName}
        type="text"
        onChange={(e) => setSuiteName(e.target.value)}
      />
      <Button
        onClick={closeModal}
        sx={{ padding: "0 !important", minWidth: "unset", float: "right" }}
      >
        <i className="material-icons">close</i>
      </Button>
      <TableContainer component={Paper} sx={{ marginTop: "2em" }}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell></TableCell>
              <TableCell>Id</TableCell>
              <TableCell>Name</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {flows.map((row) => (
              <TableRow
                key={row.id}
                sx={{
                  "&:last-child td, &:last-child th": { border: 0 },
                  cursor: "pointer",
                }}
              >
                <TableCell>
                  <Checkbox
                    onChange={(event) => handleCheckboxChange(event, row)}
                    checked={flows.findIndex((f) => f.id === row.id) !== -1}
                  />
                </TableCell>
                <TableCell component="th" scope="row" sx={{ maxWidth: "42px" }}>
                  {row.id}
                </TableCell>
                <TableCell>{row.name}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      {error && (
        <Alert sx={{ margin: "1em 0", width: "100%" }} severity="error">
          {error}
        </Alert>
      )}
      <Button
        className="add-suite-button"
        variant="contained"
        color="primary"
        sx={{ float: "right", marginTop: "12px" }}
        disabled={isLoading}
        onClick={addSuite}
      >
        Add Suite
      </Button>
    </Modal>
  );
}
