import React, { useEffect } from 'react';
import Main from 'layouts/Main';
import Box from '@mui/material/Box';
import {
  Container,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from '@mui/material';
import Typography from '@mui/material/Typography';
import axios from 'axios';
import { isArray } from 'card-validator/dist/lib/is-array';
import useToken from '../../../layouts/Main/components/Topbar/useToken';
import { getBaseUrl } from '../../../global';
import { Pass } from '../../PassPage/pass';
import Button from '@mui/material/Button';
import { SimpleDialog } from './SimpleDialog';

const History = (): JSX.Element => {
  const { token } = useToken();
  const baseUrl = getBaseUrl();

  const [openDialog, setOpenDialog] = React.useState<boolean>(false);
  const [currentPayload, setCurrentPayload] = React.useState<string>('');

  // fetch history from API
  const [history, setHistory] = React.useState<Pass[]>(null);
  useEffect(() => {
    if (!history) {
      axios
        .get(baseUrl + '/api/v1/history', {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        .then(res => {
          if (isArray(res.data)) {
            setHistory(res.data);
          }
        })
        .catch(err => {
          console.log(err);
        });
    }
  }, [history]);

  return (
    <Main>
      <Box marginTop={4} marginBottom={4}>
        <Container sx={{ maxWidth: 1500, paddingY: '0 !important' }}>
          <Box marginBottom={2}>
            <Typography variant={'h4'} fontWeight={700}>
              History
            </Typography>
            <Typography>
              Here you can find the history of your previously generated Wallet
              Passes.
            </Typography>
          </Box>
        </Container>
        <Container>
          <TableContainer>
            <Table sx={{ minWidth: 750 }} aria-label="simple table">
              <TableHead sx={{ bgcolor: 'alternate.dark' }}>
                <TableRow>
                  <TableCell>
                    <Typography
                      variant={'caption'}
                      fontWeight={700}
                      sx={{ textTransform: 'uppercase' }}
                    >
                      Id
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography
                      variant={'caption'}
                      fontWeight={700}
                      sx={{ textTransform: 'uppercase' }}
                    >
                      Generated At
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography
                      variant={'caption'}
                      fontWeight={700}
                      sx={{ textTransform: 'uppercase' }}
                    >
                      Actions
                    </Typography>
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {(!history || history.length <= 0) && (
                  <TableRow
                    sx={{
                      '&:last-child td, &:last-child th': { border: 0 }
                    }}
                  >
                    <TableCell colSpan={4}>
                      <Typography
                        color={'text.secondary'}
                        variant={'subtitle2'}
                      >
                        No history yet. Create your first Wallet Pass now.
                      </Typography>
                    </TableCell>
                  </TableRow>
                )}
                {!!history &&
                  history.map((item, i) => (
                    <TableRow
                      key={i}
                      sx={{
                        '&:last-child td, &:last-child th': { border: 0 }
                      }}
                    >
                      <TableCell>
                        <Typography
                          color={'text.secondary'}
                          variant={'subtitle2'}
                        >
                          {item.id}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography
                          color={'text.secondary'}
                          variant={'subtitle2'}
                        >
                          {new Date(
                            item.created_at * 1000
                          ).toLocaleDateString()}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Button
                          variant={'text'}
                          onClick={() => {
                            setOpenDialog(true);
                            setCurrentPayload(item.payload);
                          }}
                        >
                          See Payload
                        </Button>
                        <Button
                          variant={'text'}
                          sx={{ cursor: 'pointer', marginLeft: 2 }}
                          href={`/pass/${item.id}`}
                        >
                          See Pass
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Container>
      </Box>
      <SimpleDialog
        content={currentPayload}
        open={openDialog}
        onClose={() => setOpenDialog(false)}
      />
    </Main>
  );
};

export default History;
