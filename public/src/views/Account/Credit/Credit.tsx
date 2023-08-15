import React, { useEffect, useState } from 'react';
import Main from 'layouts/Main';
import Box from '@mui/material/Box';
import Container from '../../../components/Container';
import Grid from '@mui/material/Grid';
import { Card, Stack } from '@mui/material';
import Typography from '@mui/material/Typography';
import useToken, {
  Customer
} from '../../../layouts/Main/components/Topbar/useToken';
import axios from 'axios';
import { CreditCard } from '@mui/icons-material';
import Button from '@mui/material/Button';
import { loadStripe } from '@stripe/stripe-js/pure';
import { useNavigate } from 'react-router-dom';
import ConfettiExplosion from 'react-confetti-explosion';
import { getBaseUrl, useQuery } from '../../../global';

const Credit = (): JSX.Element => {
  const stripePromise = loadStripe(
    'pk_live_51MMpuYG7ggTrzU4v6lWetJGeHqGBZ2XBxKRu5xCpJCXSorsbid2ymK9yw6IEx6cGPJfK0P7zgjzTf0NrYebK6Y0p00OJFPRFP7'
  );
  const navigate = useNavigate();
  const { token, setToken } = useToken();
  const [profile, setProfile] = useState<Customer>(null);
  const [checkoutSessionId, setCheckoutSessionId] = useState<string>('');
  const query = useQuery();
  const isSuccess = query.get('success') === 'true';
  const baseUrl = getBaseUrl();

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
          console.error(err);
          navigate('/');
          setToken(null);
        });
    }
  }, [token]);

  useEffect(() => {
    if (profile) {
      axios
        .get(baseUrl + '/api/v1/checkout', {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        .then(res => {
          setCheckoutSessionId(res.data?.sessionId ?? '');
        })
        .catch(err => {
          console.error(err);
          navigate('/error');
        });
    }
  }, [profile]);

  return (
    <Main>
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
              {isSuccess && (
                <Box>
                  <ConfettiExplosion />
                  <Typography
                    color={'text.primary'}
                    variant={'subtitle2'}
                    textAlign={'center'}
                  >
                    <b>Thank you for your purchase!</b> We added the credits to
                    your account.
                  </Typography>
                </Box>
              )}
              {!isSuccess && (
                <Box>
                  <Box>
                    <Typography
                      color={'text.primary'}
                      variant={'subtitle2'}
                      textAlign={'center'}
                    >
                      You have <b>{profile?.credit_tokens ?? '-'} credits</b>.
                      Click on the button bellow to add 1000x credits for $1.99.
                    </Typography>
                  </Box>
                  <Box flexGrow={1} />
                  <Stack
                    spacing={2}
                    marginTop={4}
                    width={1}
                    alignItems={'center'}
                  >
                    <Box
                      display={'flex'}
                      justifyContent={'center'}
                      alignItems={'center'}
                    >
                      {!!checkoutSessionId && (
                        <Button
                          variant="contained"
                          color="primary"
                          onClick={() => {
                            void stripePromise.then(stripe => {
                              void stripe.redirectToCheckout({
                                sessionId: checkoutSessionId
                              });
                            });
                          }}
                          startIcon={<CreditCard />}
                        >
                          Add Credits
                        </Button>
                      )}
                    </Box>
                  </Stack>
                </Box>
              )}
            </Card>
          </Grid>
        </Container>
      </Box>
    </Main>
  );
};

export default Credit;
