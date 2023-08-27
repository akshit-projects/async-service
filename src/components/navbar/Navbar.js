import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Link, useNavigate } from 'react-router-dom';
import { AppBar, Toolbar, Typography, Button } from '@material-ui/core';
import constants from '../../constants/constants';
import { checkLoginState, logout } from '../auth/auth-utils';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    'max-width': '80%',
    margin: 'auto',
    height: '68px'
  },
  title: {
    flexGrow: 1,
    'fontSize': '20px',
    fontWeight: 500,
  },
  linkButton: {
    margin: theme.spacing(0),
    padding: theme.spacing(1, 2),
    'fontSize': '12px',
    fontFamily: 'helvetica',
  },
}));

function Navbar() {
  const classes = useStyles();
  const navigate = useNavigate();

  const cleanAndLogout = () => {
    navigate('/login');
    logout();
  }

  return (
    <div className={classes.root}>
      <AppBar position="static" style={{ backgroundColor: '#ffffff', color: '#000000', boxShadow: 'none' }}>
        <Toolbar>
          <Typography className={classes.title}>
            Akshit Helper
          </Typography>
          <Button color="inherit" className={classes.linkButton}>
            <Link to={constants.PATHS.FLOWS}>Flow</Link>
          </Button>
          {checkLoginState() && <Button color="inherit" onClick={cleanAndLogout} className={classes.linkButton}>
            Logout
          </Button>}
        </Toolbar>
      </AppBar>
    </div>
  );
}

export default Navbar;
