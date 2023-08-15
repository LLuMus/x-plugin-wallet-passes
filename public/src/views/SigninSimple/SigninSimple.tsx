import React, { useEffect, useState } from 'react';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Main from 'layouts/Main';
import Container from 'components/Container';
import { getBaseUrl, useQuery } from '../../global';
import { useNavigate } from 'react-router-dom';
import { Card } from '@mui/material';
import Typography from '@mui/material/Typography';
import { GoogleLogin } from '@react-oauth/google';
import useToken from '../../layouts/Main/components/Topbar/useToken';
import axios from 'axios';

const SigninSimple = (): JSX.Element => {
  const navigate = useNavigate();
  const query = useQuery();
  const redirectUri = query.get('redirect_uri');
  const state = query.get('state');
  if (!redirectUri || !state) {
    navigate('/');
  }

  const { token, setToken } = useToken();
  const [googleToken, setGoogleToken] = useState('');
  const baseUrl = getBaseUrl();

  const responseMessage = response => {
    setGoogleToken(response.credential);
  };

  useEffect(() => {
    if (token) {
      window.location.href = redirectUri + '?code=' + token + '&state=' + state;
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
    <Main showSignin={false}>
      <Box bgcolor={'primary.grey'}>
        <Container>
          <Grid item xs={12} sm={4} display={'flex'} justifyContent={'center'}>
            <Card
              sx={{
                p: { xs: 2, md: 4 },
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                maxWidth: 360,
                width: 1,
                height: 1,
                background: 'transparent'
              }}
            >
              {token && (
                <Box>
                  <Typography
                    color={'text.primary'}
                    variant={'subtitle2'}
                    textAlign={'center'}
                  >
                    Great! Hold on a second we are going back to ChatGPT.
                  </Typography>
                </Box>
              )}

              {!token && (
                <Box>
                  <Typography
                    color={'text.primary'}
                    variant={'subtitle2'}
                    textAlign={'center'}
                  >
                    Welcome! Please Sign in with Google to continue to ChatGPT.
                  </Typography>
                </Box>
              )}
              {!token && (
                <Box
                  display={'flex'}
                  justifyContent={'center'}
                  alignItems={'center'}
                  marginTop={2}
                >
                  <GoogleLogin onSuccess={responseMessage} />
                </Box>
              )}
            </Card>
          </Grid>
        </Container>
      </Box>
    </Main>
  );
};

export default SigninSimple;
