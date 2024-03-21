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
  TablePagination,
  TableFooter,
} from "@mui/material";
import constants from "../../constants/constants";
import CheckBox from "@mui/material/Checkbox";
import AddSuiteModal from "../test-suite/AddSuiteModal";

let timeoutId;
function Flows() {
  const [data, setData] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const [page, setPage] = useState(0);
  const [selectedFlows, setSelectedFlows] = useState([]);
  const [suiteModal, setSuiteModal] = useState(false);

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

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleRowClick = (id) => {
    navigate(`${constants.PATHS.FLOWS}/${id}`);
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

  const addTestSuite = (e) => {
    setSuiteModal(true);
    console.log(selectedFlows);
  };

  const handleCheckboxChange = (event, flow) => {
    if (event.stopPropagation) {
      event.stopPropagation();
    } else {
      event.cancelBubble = true;
    }
    if (event.target.checked) {
      setSelectedFlows([...selectedFlows, flow]);
    } else {
      setSelectedFlows(selectedFlows.filter((f) => f.id !== flow.id));
    }
  };

  const openAddFlow = (event) => {
    navigate(constants.PATHS.ADD_FLOW);
  }

  return (
    <Container>
        
      <Grid container spacing={2} alignItems="center" marginTop={2}>
        <Grid item xs={1.5}>
            <h1>Flows</h1>
        </Grid>
        <Grid item xs={4}>
          <TextField
            label="Search"
            variant="outlined"
            fullWidth
            onChange={(e) => filterResults(e.target.value)}
          />
        </Grid>
        {selectedFlows.length ? ( // todo move it out of here
          <Grid xs={1} sx={{ verticalAlign: "middle" }}>
            <Button variant="contained" onClick={addTestSuite}>
              Create Test Suite
            </Button>
          </Grid>
        ) : (
          <></>
        )}
        <Grid
          item
          xs={6.5}
          style={{ textAlign: "right" }}
        >
          {/* <Fab // TODO fix this
            color="primary"
            aria-label="Add"
            component={Link}
            to={constants.PATHS.ADD_FLOW}
          >
            <i className="material-icons fs-4">add</i>
          </Fab> */}
          <Button variant="contained" onClick={openAddFlow}>Add Flow</Button>
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
              <TableCell width="15%">Steps</TableCell>
              <TableCell width="20%">Created At</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map((row) => (
              <TableRow
                key={row.id}
                sx={{
                  cursor: "pointer",
                }}
              >
                <TableCell width="10%">
                  <CheckBox
                    onChange={(event) => handleCheckboxChange(event, row)}
                    checked={
                      selectedFlows.findIndex((f) => f.id === row.id) !== -1
                    }
                  />
                </TableCell>
                <TableCell
                  width="20%"
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
                <TableCell width="20%" onClick={() => handleRowClick(row.id)}>
                  {new Date(parseInt(row.createdAt) * 1000).toLocaleString()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
          <TableFooter>
            <TableRow>
              <TablePagination
                align="right"
                rowsPerPageOptions={[]}
                count={data.length}
                rowsPerPage={6}
                page={page}
                SelectProps={{
                  inputProps: {
                    "aria-label": "rows per page",
                  },
                  native: false,
                }}
                onPageChange={handleChangePage}
              />
            </TableRow>
          </TableFooter>
        </Table>
      </TableContainer>
      <AddSuiteModal
        isOpen={suiteModal}
        flows={selectedFlows}
        closeModal={() => setSuiteModal(false)}
        handleCheckboxChange={handleCheckboxChange}
      />
    </Container>
  );
}

export default Flows;
