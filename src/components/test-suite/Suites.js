import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import {
  Container,
  Table,
  Grid,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  TableHead,
  Paper,
  TextField,
  Fab,
  Alert,
  Button,
} from "@mui/material";
import constants from "../../constants/constants";
import SuiteRow from "./SuiteRow";

let timeoutId;
function Suites() {
  const [data, setData] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const [selectedFlows, setSelectedFlows] = useState([]);
  const [suiteModal, setSuiteModal] = useState(false);

  useEffect(() => {
    const getFlows = async () => {
      const opts = {
        url: `${constants.BACKEND_URL}/api/v1/suite`,
        method: "GET",
      };
      if (searchQuery) {
        opts.url = `${opts.url}?search=${searchQuery}`;
      }
      await axios(opts)
        .then((response) => {
          setData(response.data);
          setError();
        })
        .catch((err) => {
          if (
            err.response?.status === constants.HTTP_STATUS_RATE_LIMIT_EXCEEDED
          ) {
            setError("Please wait before retrying");
          } else {
            setError(err.message);
          }
        });
    };
    getFlows();
  }, [searchQuery]);

  const handleRowClick = (id) => {
    navigate(`/suite/${id}`);
  };

  const filterResults = (query) => {
    if (timeoutId) {
      clearTimeout(timeoutId);
      timeoutId = null;
    }
    timeoutId = setTimeout(() => {
      setSearchQuery(query);
    }, 1000);
  };

  return (
    <Container>
      <Grid container spacing={1} alignItems="center">
        <Grid item xs={4}>
          <TextField
            label="Search"
            variant="outlined"
            fullWidth
            onChange={(e) => filterResults(e.target.value)}
          />
        </Grid>
      </Grid>
      {error && (
        <Alert sx={{ margin: "1em 0" }} severity="error">
          {error}
        </Alert>
      )}
      <TableContainer component={Paper} sx={{ marginTop: "2em" }}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell width="10%"></TableCell>
              <TableCell width="10%">Id</TableCell>
              <TableCell>Name</TableCell>
              <TableCell width="15%">Flows</TableCell>
              <TableCell width="20%">Created At</TableCell>
              <TableCell  width="5%" align="right"></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map((row) => (
                <SuiteRow row={row} />
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
}

export default Suites;
