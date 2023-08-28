import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
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
import CheckBox from "@mui/material/Checkbox";

let timeoutId;
function Flows() {
  const [data, setData] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const [selectedRows, setSelectedRows] = useState([]);

  useEffect(() => {
    const getFlows = async () => {
      const opts = {
        url: `${constants.BACKEND_URL}/api/v1/flow`,
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
    navigate(`${constants.PATHS.API_PREFIX}/${id}`);
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

  const handleCheckboxChange = (event, id) => {
    if (event.stopPropagation) {
      event.stopPropagation();
    } else {
      event.cancelBubble = true;
    }
    if (event.target.checked) {
      setSelectedRows([...selectedRows, id]);
    } else {
      setSelectedRows(selectedRows.filter((rowId) => rowId !== id));
    }
  };

  return (
    <Container>
      <Grid container spacing={2}>
        <Grid item xs={4}>
          <TextField
            label="Search"
            variant="outlined"
            fullWidth
            onChange={(e) => filterResults(e.target.value)}
          />
        </Grid>
        <Grid item xs={8} style={{ textAlign: "right" }}>
          <Fab color="primary" aria-label="Add" component={Link} to={constants.PATHS.ADD_FLOW}>
            <i className="material-icons fs-4">add</i>
          </Fab>
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
            <TableCell></TableCell>
              <TableCell>Id</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Steps</TableCell>
              <TableCell>Created At</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map((row) => (
              <TableRow
                key={row.id}
                sx={{
                  "&:last-child td, &:last-child th": { border: 0 },
                  cursor: "pointer",
                }}
              >
                <TableCell>
                  <CheckBox
                    onChange={(event) => handleCheckboxChange(event, row.id)}
                    checked={selectedRows.includes(row.id)}
                  />
                </TableCell>
                <TableCell
                  component="th"
                  scope="row"
                  onClick={() => handleRowClick(row.id)}
                  sx={{ maxWidth: "42px" }}
                >
                  {row.id}
                </TableCell>
                <TableCell onClick={() => handleRowClick(row.id)}>
                  {row.name}
                </TableCell>
                <TableCell onClick={() => handleRowClick(row.id)}>
                  {row.steps.length}
                </TableCell>
                <TableCell onClick={() => handleRowClick(row.id)}>
                  {new Date(parseInt(row.createdAt) * 1000).toLocaleString()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
}

export default Flows;
