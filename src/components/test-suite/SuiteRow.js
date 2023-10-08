import React, { useState } from "react";
import PropTypes from "prop-types";
import Box from "@mui/material/Box";
import Collapse from "@mui/material/Collapse";
import IconButton from "@mui/material/IconButton";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Typography from "@mui/material/Typography";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import constants from "../../constants/constants";
import axios from "axios";
import { Alert } from "@mui/material";
import "./Suite.css";

function SuiteRow(props) {
  const { row } = props;
  const [open, setOpen] = React.useState(false);
  const [flowData, setFlowData] = useState([]);
  const [error, setError] = useState(null);
  const runSuite = () => {
    console.log(row);
  }
  const openRow = (isOpen) => {
    setOpen(isOpen);
    console.log(flowData);
    if (flowData && flowData.length) {
      return;
    }

    const options = {
      url: `${constants.BACKEND_URL}/api/v1/flow?ids=${row.flowIds.join(",")}`,
      method: "GET",
    };

    axios(options)
      .then((data) => {
        setError();
        setFlowData(data.data);
      })
      .catch((err) => {
        setError("Unable to get flows data.");
      });
  };

  return (
    <React.Fragment>
      <TableRow sx={{ "& > *": { borderBottom: "unset" } }}>
        <TableCell>
          <IconButton
            aria-label="expand row"
            size="small"
            onClick={() => openRow(!open)}
          >
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell>{row.id}</TableCell>
        <TableCell width="30%">{row.name}</TableCell>
        <TableCell>{row.flowIds.length}</TableCell>
        <TableCell>
          {new Date(parseInt(row.createdAt) * 1000).toLocaleString()}
        </TableCell>
        <TableCell width="5%">
          <i onClick={runSuite} className="material-icons run-suite">
            play_arrow
          </i>
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              {error ? (
                <Alert sx={{ margin: "1em 0" }} severity="error">
                  {error}
                </Alert>
              ) : (
                <>
                  <Typography variant="h6" gutterBottom component="div">
                    Flows
                  </Typography>
                  <Table size="small" aria-label="purchases">
                    <TableHead>
                      <TableRow>
                        <TableCell>Id</TableCell>
                        <TableCell>Name</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {flowData?.map((flow) => (
                        <TableRow key={flow.date}>
                          <TableCell>{flow.id}</TableCell>
                          <TableCell>{flow.name}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </>
              )}
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
}

SuiteRow.propTypes = {
  row: PropTypes.shape({
    calories: PropTypes.number.isRequired,
    carbs: PropTypes.number.isRequired,
    fat: PropTypes.number.isRequired,
    history: PropTypes.arrayOf(
      PropTypes.shape({
        amount: PropTypes.number.isRequired,
        customerId: PropTypes.string.isRequired,
        date: PropTypes.string.isRequired,
      })
    ).isRequired,
    name: PropTypes.string.isRequired,
    price: PropTypes.number.isRequired,
    protein: PropTypes.number.isRequired,
  }).isRequired,
};

export default SuiteRow;
