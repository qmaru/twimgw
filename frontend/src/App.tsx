import { useState } from 'react'

import Container from '@mui/material/Container'
import Box from '@mui/material/Box'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import FormControl from '@mui/material/FormControl'
import FormControlLabel from '@mui/material/FormControlLabel'
import Checkbox from '@mui/material/Checkbox'
import Stack from '@mui/material/Stack'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'

import { TwitterCore } from "../wailsjs/go/apps/App"

function App() {
  const [username, setUsername] = useState<string>("")
  const [exclude, setExclude] = useState<boolean>(true)
  const [maxResults, setMaxResults] = useState<number>(100)
  const [startID, setStartID] = useState<string>("")
  const [storagePath, setStoragePath] = useState<string>("")
  const [socks5, setSocks5] = useState<string>("")

  const [usernameError, setUsernameError] = useState<boolean>(false)
  const [usernameLabel, setUsernameLabel] = useState<string>("username")

  const [maxResultsError, setMaxResultsError] = useState<boolean>(false)
  const [maxResultsLabel, setMaxResultsLabel] = useState<string>("maxResults")

  const [loading, setLoading] = useState<boolean>(false)
  const [coreResults, setCoreResults] = useState<string>("")

  const UsernameChange = (event: any) => {
    setUsername(event.target.value)
  }

  const ExcludeChange = (event: any) => {
    setExclude(event.target.checked)
  }

  const MaxResultsChange = (event: any) => {
    setMaxResults(event.target.value)
  }

  const StartIDChange = (event: any) => {
    setStartID(event.target.value)
  }

  const StoragePathChange = (event: any) => {
    setStoragePath(event.target.value)
  }

  const Socks5Change = (event: any) => {
    setSocks5(event.target.value)
  }

  const ResetInput = () => {
    setUsername("")
    setExclude(true)
    setMaxResults(100)
    setStartID("")
    setStoragePath("")
    setSocks5("")
  }

  const MediaDownload = () => {
    if (username.trim() === "") {
      setUsernameError(true)
      setUsernameLabel("username is empty")
      return false
    } else {
      setUsernameError(false)
      setUsernameLabel("username")
    }

    if (maxResults < 5 || maxResults > 100) {
      setMaxResultsError(true)
      setMaxResultsLabel("max results must in (5, 100)")
      return false
    } else {
      setMaxResultsError(false)
      setMaxResultsLabel("maxResults")
    }

    setLoading(true)

    const Username = username.trim()
    const MaxResults = String(maxResults)
    const StartID = startID.trim()
    const StoragePath = storagePath.trim()
    const Socks5 = socks5.trim()

    const body: any = {
      "username": Username,
      "max_results": MaxResults,
      "start_id": StartID,
      "exclude": exclude,
      "socks5": Socks5,
      "storage_path": StoragePath,
    }

    setCoreResults("Downloading...")
    TwitterCore(body)
      .then((res) => {
        setLoading(false)
        const status = res.status
        const message = res.message
        if (status === 1) {
          const output = res.data
          setCoreResults("Finished: " + output)
        } else {
          setCoreResults("Downloading error: " + message)
        }
      })
  }

  return (
    <Container key="App-Wrapper"
      sx={{
        width: 400,
        paddingTop: 2
      }}
    >
      <Container key="App-Form">
        <Box>
          <Stack
            direction="column"
            justifyContent="center"
            alignItems="center"
            spacing={2}
          >
            <Button
              variant="outlined"
              color="primary"
              size="small"
              onClick={() => ResetInput()}
              disabled={loading}
            >
              reset
            </Button>
            <TextField
              id="username"
              fullWidth
              label={usernameLabel}
              error={usernameError}
              value={username}
              onChange={event => UsernameChange(event)}
              margin="normal"
              variant="standard"
            />
            <TextField
              id="maxResults"
              fullWidth
              label={maxResultsLabel}
              error={maxResultsError}
              placeholder="100"
              value={maxResults}
              onChange={event => MaxResultsChange(event)}
              margin="normal"
              variant="standard"
            />
            <TextField
              id="startID"
              fullWidth
              label="startID"
              value={startID}
              onChange={event => StartIDChange(event)}
              margin="normal"
              variant="standard"
            />
            <TextField
              id="storagePath"
              fullWidth
              label="storagePath"
              value={storagePath}
              onChange={event => StoragePathChange(event)}
              margin="normal"
              variant="standard"
            />
            <TextField
              id="socks5"
              fullWidth
              label="socks5"
              value={socks5}
              onChange={event => Socks5Change(event)}
              margin="normal"
              variant="standard"
            />
            <FormControl fullWidth>
              <FormControlLabel
                label="exclude"
                control={
                  <Checkbox
                    checked={exclude}
                    onChange={ExcludeChange}
                  />
                }
              />
            </FormControl>
          </Stack>
        </Box>

        <Box sx={{ paddingTop: 2 }}>
          <Button
            fullWidth
            variant="contained"
            color="primary"
            size="large"
            onClick={() => MediaDownload()}
            disabled={loading}
          >
            Download
          </Button>

          <Box sx={{ paddingTop: 2 }} hidden={coreResults === ""}>
            <Paper elevation={1} sx={{ padding: 2 }}>
              <Typography sx={{ wordWrap: 'break-word' }}>
                {coreResults}
              </Typography>
            </Paper>
          </Box>
        </Box>
      </Container>
    </Container>
  )
}

export default App
