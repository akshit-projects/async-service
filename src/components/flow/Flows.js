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
} from "@mui/material";
import constants from "../../constants/constants";

function Flows() {
  const [data, setData] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const getFlows = async () => {
      const opts = {
        url: `${constants.BACKEND_URL}/api/v1/flow`,
        method: "GET",
      };
      await axios(opts)
        .then((response) => {
          setData(response.data);
          setError();
        })
        .catch((err) => {
          if (err.response?.status === constants.HTTP_STATUS_RATE_LIMIT_EXCEEDED) {
            setError("Please wait before retrying");
          } else {
            setError(err.message);
          }
        });
    };
    getFlows();
  }, []);

  const filteredData = data.filter((item) =>
    item.name.toLowerCase().includes(searchQuery.toLowerCase())
  );
  const handleRowClick = (id) => {
    navigate(`${constants.PATHS.API_PREFIX}/${id}`);
  };

  return (
    <Container>
      <Grid container spacing={2}>
        <Grid item xs={4}>
          <TextField
            label="Search"
            variant="outlined"
            fullWidth
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </Grid>
        <Grid item xs={8} style={{ textAlign: "right" }}>
          <Fab color="primary" aria-label="Add" component={Link} to={constants.PATHS.FLOWS}>
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
              <TableCell>Id</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Steps</TableCell>
              <TableCell>Created At</TableCell>
              <TableCell>Creator</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredData.map((row) => (
              <TableRow
                onClick={() => handleRowClick(row.id)}
                key={row.id}
                sx={{ "&:last-child td, &:last-child th": { border: 0 }, cursor: 'pointer' }}
              >
                <TableCell component="th" scope="row" sx={{ maxWidth: "42px" }}>
                  {row.id}
                </TableCell>
                <TableCell>{row.name}</TableCell>
                <TableCell>{row.steps.length}</TableCell>
                <TableCell>{row.createdAt}</TableCell>
                <TableCell>{row.creator}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
}

export default Flows;
