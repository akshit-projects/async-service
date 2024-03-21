import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import { Link, useNavigate } from "react-router-dom";
import { AppBar, Toolbar, Typography, Button } from "@material-ui/core";
import constants from "../../constants/constants";
import { checkLoginState, logout } from "../auth/auth-utils";

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    "max-width": "100%",
    margin: "auto",
    height: "68px",
  },
  title: {
    flexGrow: 1,
    fontSize: "20px",
    fontWeight: 500,
  },
  linkButton: {
    margin: theme.spacing(0),
    padding: theme.spacing(1, 2),
    fontSize: "12px",
    fontFamily: "helvetica",
  },
  innerBar: {
    width: '80%',
    margin: 'auto'
  },
}));

function Navbar() {
  const classes = useStyles();
  const navigate = useNavigate();

  const cleanAndLogout = () => {
    navigate("/login");
    logout();
  };

  window.addEventListener("userLogout", () => {
    navigate("/login");
  });

  return (
    <div className={classes.root}>
      <AppBar
        position="static"
        style={{
          backgroundColor: "#ffffff",
          color: "#000000",
          boxShadow: "0 4px 2px -2px rgba(0,0,0,.2);",
        }}
      >
        <Toolbar className={classes.innerBar}>
          <Typography className={classes.title}>Akshit Helper</Typography>
          {checkLoginState() && (
            <>
              <Button color="inherit" className={classes.linkButton}>
                <Link to={constants.PATHS.FLOWS}>Flows</Link>
              </Button>
              <Button color="inherit" className={classes.linkButton}>
                <Link to={constants.PATHS.SUITES}>Suites</Link>
              </Button>
              <Button
                color="inherit"
                onClick={cleanAndLogout}
                className={classes.linkButton}
              >
                Logout
              </Button>
            </>
          )}
        </Toolbar>
      </AppBar>
    </div>
  );
}

export default Navbar;
