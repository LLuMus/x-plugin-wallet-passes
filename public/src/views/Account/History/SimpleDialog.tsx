import { Dialog, DialogTitle } from '@mui/material';
import Box from '@mui/material/Box';

export interface SimpleDialogProps {
  open: boolean;
  content: string;
  onClose: (value: string) => void;
}

export function SimpleDialog(props: SimpleDialogProps) {
  const { onClose, open } = props;

  const handleClose = () => {
    onClose(null);
  };

  const prepared = !!props.content
    ? JSON.stringify(JSON.parse(props.content), null, 2)
    : null;

  return (
    <Dialog onClose={handleClose} open={open}>
      <DialogTitle>Payload Transmitted</DialogTitle>
      <Box>
        <pre
          style={{
            width: '100%',
            display: 'block',
            whiteSpace: 'break-spaces',
            paddingLeft: '24px',
            paddingRight: '24px'
          }}
        >
          {prepared}
        </pre>
      </Box>
    </Dialog>
  );
}
