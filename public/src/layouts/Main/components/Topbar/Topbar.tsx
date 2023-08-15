import React, { useEffect, useState } from 'react';
import Box from '@mui/material/Box';
import { useTheme } from '@mui/material/styles';
import { GoogleLogin } from '@react-oauth/google';
import useToken, { Customer } from './useToken';
import axios from 'axios';
import Button from '@mui/material/Button';
import { Badge } from '@mui/material';
import { CurrencyExchangeTwoTone, History } from '@mui/icons-material';
import { getBaseUrl } from '../../../../global';

interface Props {
  showSignin?: boolean;
}

const Topbar = ({ showSignin = true }: Props): JSX.Element => {
  const theme = useTheme();
  const { mode } = theme.palette;
  const { token, setToken } = useToken();
  const [googleToken, setGoogleToken] = useState('');
  const [profile, setProfile] = useState<Customer>(null);
  const baseUrl = getBaseUrl();

  const responseMessage = response => {
    setGoogleToken(response.credential);
  };

  useEffect(() => {
    if (token) {
      axios
        .get(baseUrl + '/api/v1/me', {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        .then(res => {
          setProfile(res.data);
        })
        .catch(err => {
          console.log(err);
          setToken(null);
        });
    }
  }, [token]);

  useEffect(() => {
    if (googleToken && !token) {
      axios
        .post(baseUrl + '/auth/google', {
          access_token: googleToken
        })
        .then(res => {
          setToken(res.data.auth_code);
        })
        .catch(err => console.log(err));
    }
  }, [googleToken, token]);

  return (
    <Box
      display={'flex'}
      justifyContent={showSignin ? 'space-between' : 'center'}
      alignItems={'center'}
      width={1}
    >
      <Box
        display={'flex'}
        component="a"
        href="/"
        title="theFront"
        width={{ xs: 100, md: 120 }}
      >
        <Box
          component={'img'}
          src={
            mode === 'light'
              ? 'https://storage.googleapis.com/walletpasses/logo_new.svg'
              : 'https://storage.googleapis.com/walletpasses/logo_new_inverted.svg'
          }
          height={1}
          width={1}
        />
      </Box>
      {showSignin && (
        <Box marginLeft={4}>
          {!profile && <GoogleLogin onSuccess={responseMessage} />}
          {profile && (
            <Box>
              <Button
                variant="text"
                color="primary"
                component="a"
                target="_self"
                href="/history"
                sx={{
                  marginRight: 2
                }}
                startIcon={<History />}
              >
                History
              </Button>{' '}
              <Button
                variant="text"
                color="primary"
                component="a"
                target="_self"
                href="/credit"
                startIcon={<CurrencyExchangeTwoTone />}
              >
                <Badge badgeContent={profile.credit_tokens} color="secondary">
                  Credits
                </Badge>
              </Button>
            </Box>
          )}
        </Box>
      )}
    </Box>
  );
};

export default Topbar;
