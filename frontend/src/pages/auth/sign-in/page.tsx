import React, { useState } from 'react';
import { useFormik } from 'formik';
import * as yup from 'yup';
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Container,
  Divider,
  FormControl,
  FormLabel,
  IconButton,
  Input,
  InputAdornment,
  Link,
  Paper,
  TextField,
  Typography,
} from '@mui/material';
import { Visibility, VisibilityOff, ErrorOutline } from '@mui/icons-material';
import { useAuth } from '@/contexts/AuthContext'; // Assuming this path is correct
import { useNavigate } from 'react-router-dom';

// --- Helper Components & SVGs (from the second file) ---

// A placeholder for the Logo component
const Logo = () => (
  <Typography variant="h4" component="h1" fontWeight="bold" color="primary">
    PantryOS
  </Typography>
);

const googleSVG = () => (
  <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M19.6169 10.2876C19.6169 9.60932 19.5561 8.95714 19.443 8.33105H10.4343V12.0354H15.5822C15.3561 13.2267 14.6778 14.2354 13.6604 14.9137V17.3224H16.7648C18.5735 15.6528 19.6169 13.2006 19.6169 10.2876Z" fill="#4285F4" />
    <path d="M10.4346 19.6346C13.0172 19.6346 15.1825 18.7825 16.7651 17.3216L13.6607 14.9129C12.8086 15.4868 11.7216 15.8346 10.4346 15.8346C7.94768 15.8346 5.83464 14.1564 5.07812 11.8955H1.89551V14.3651C3.46942 17.4868 6.69551 19.6346 10.4346 19.6346Z" fill="#34A853" />
    <path d="M5.07832 11.8866C4.88702 11.3127 4.77398 10.704 4.77398 10.0692C4.77398 9.4344 4.88702 8.8257 5.07832 8.25179V5.78223H1.89572C1.24354 7.06918 0.869629 8.52136 0.869629 10.0692C0.869629 11.617 1.24354 13.0692 1.89572 14.3561L4.37398 12.4257L5.07832 11.8866Z" fill="#FBBC05" />
    <path d="M10.4346 4.31358C11.8433 4.31358 13.0955 4.80054 14.0955 5.73967L16.8346 3.00054C15.1738 1.45271 13.0172 0.504883 10.4346 0.504883C6.69551 0.504883 3.46942 2.65271 1.89551 5.78314L5.07812 8.25271C5.83464 5.99184 7.94768 4.31358 10.4346 4.31358Z" fill="#EA4335" />
  </svg>
);

const githubSVG = () => (
  <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path fillRule="evenodd" clipRule="evenodd" d="M10 0.0693359C4.475 0.0693359 0 4.54434 0 10.0693C0 14.4943 2.8625 18.2318 6.8375 19.5568C7.3375 19.6443 7.525 19.3443 7.525 19.0818C7.525 18.8443 7.5125 18.0568 7.5125 17.2193C5 17.6818 4.35 16.6068 4.15 16.0443C4.0375 15.7568 3.55 14.8693 3.125 14.6318C2.775 14.4443 2.275 13.9818 3.1125 13.9693C3.9 13.9568 4.4625 14.6943 4.65 14.9943C5.55 16.5068 6.9875 16.0818 7.5625 15.8193C7.65 15.1693 7.9125 14.7318 8.2 14.4818C5.975 14.2318 3.65 13.3693 3.65 9.54434C3.65 8.45684 4.0375 7.55684 4.675 6.85684C4.575 6.60684 4.225 5.58184 4.775 4.20684C4.775 4.20684 5.6125 3.94434 7.525 5.23184C8.325 5.00684 9.175 4.89434 10.025 4.89434C10.875 4.89434 11.725 5.00684 12.525 5.23184C14.4375 3.93184 15.275 4.20684 15.275 4.20684C15.825 5.58184 15.475 6.60684 15.375 6.85684C16.0125 7.55684 16.4 8.44434 16.4 9.54434C16.4 13.3818 14.0625 14.2318 11.8375 14.4818C12.2 14.7943 12.5125 15.3943 12.5125 16.3318C12.5125 17.6693 12.5 18.7443 12.5 19.0818C12.5 19.3443 12.6875 19.6568 13.1875 19.5568C17.1375 18.2318 20 14.4818 20 10.0693C20 4.54434 15.525 0.0693359 10 0.0693359Z" fill="#1B1F23" />
  </svg>
);


// --- Merged Login Component ---

const validationSchema = yup.object({
  email: yup.string().required('Email is required').email('Enter a valid email'),
  password: yup.string().required('Password is required'),
});

const Login: React.FC = () => {
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const { login } = useAuth();
  const navigate = useNavigate();

  const formik = useFormik({
    initialValues: {
      email: 'info@pantryOS.dev',
      password: 'pantryOS',
    },
    validationSchema,
    onSubmit: async (values) => {
      
      setError('');
      setIsLoading(true);
      try {
        await login(values.email, values.password);
        navigate('/dashboard'); // Redirect on success
      } catch (err: any) {
        setError(err.response?.data?.error || 'Login failed. Please check your credentials.');
      } finally {
        setIsLoading(false);
      }
    },
  });

  const handleRegisterClick = () => {
    navigate('/auth/register');
  };

  return (
    <Container component="main" maxWidth="sm">
      <Box sx={{ marginTop: 8, display: 'flex', justifyContent: 'center' }}>
        <Paper
          elevation={4}
          sx={{
            padding: { xs: 3, sm: 5 },
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            width: '100%',
            maxWidth: '480px',
            borderRadius: '16px',
          }}
        >
          <Box sx={{ mb: 4 }}>
            <Logo />
          </Box>

          <Typography component="h1" variant="h5" sx={{ mb: 1 }}>
            Sign in
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            Access your account quickly and securely.
          </Typography>

          {error && (
            <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={formik.handleSubmit} sx={{ width: '100%' }}>
            <TextField
              fullWidth
              id="email"
              name="email"
              label="Email Address"
              margin="normal"
              autoComplete="email"
              autoFocus
              value={formik.values.email}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.email && Boolean(formik.errors.email)}
              helperText={formik.touched.email && formik.errors.email}
              disabled={isLoading}
            />
            <TextField
              fullWidth
              id="password"
              name="password"
              label="Password"
              type={showPassword ? 'text' : 'password'}
              margin="normal"
              autoComplete="current-password"
              value={formik.values.password}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.password && Boolean(formik.errors.password)}
              helperText={formik.touched.password && formik.errors.password}
              disabled={isLoading}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      aria-label="toggle password visibility"
                      onClick={() => setShowPassword(!showPassword)}
                      onMouseDown={(e) => e.preventDefault()}
                      edge="end"
                    >
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
            
            <Box sx={{ textAlign: 'right', width: '100%', my: 1 }}>
                <Link href="/auth/password-reset" variant="body2">
                    Forgot password?
                </Link>
            </Box>

            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 2, mb: 2, height: 48 }}
              disabled={isLoading}
            >
              {isLoading ? <CircularProgress size={24} color="inherit" /> : 'Sign In'}
            </Button>

            <Typography variant="body2" color="text.secondary" align="center" sx={{my: 2}}>
                By clicking continue, you agree to our{' '}
                <Link href="/auth/terms-and-conditions" underline="always">
                    Terms of Service
                </Link>{' '}
                and{' '}
                <Link href="/auth/privacy-policy" underline="always">
                    Privacy Policy
                </Link>
                .
            </Typography>

            <Divider sx={{ width: '100%', my: 2 }} />
          </Box>
                      <Box sx={{ textAlign: 'center' }}>
              <Typography variant="body2" color="text.secondary">
                Don't have an account?{' '}
                <Link component="button" variant="body2" onClick={handleRegisterClick} disabled={isLoading}>
                  Sign Up
                </Link>
              </Typography>
            </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default Login;
