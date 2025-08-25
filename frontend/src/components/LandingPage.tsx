import React from 'react';
import {
  Box,
  Button,
  Container,
  Typography,
  Card,
  CardContent,
  AppBar,
  Toolbar,
  useTheme,
  useMediaQuery,
  Paper,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Chip,
  Stack,
  Divider,
} from '@mui/material';
import {
  Inventory,
  Analytics,
  Security,
  Speed,
  CloudSync,
  Dashboard as DashboardIcon,
  Login,
  PersonAdd,
  TrendingUp,
  Notifications,
  Assessment,
  Settings,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

const LandingPage: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const navigate = useNavigate();

  const features = [
    {
      icon: <Inventory sx={{ fontSize: 40, color: theme.palette.primary.main }} />,
      title: 'Inventory Management',
      description: 'Track your stock levels, manage products, and monitor inventory in real-time.',
    },
    {
      icon: <Analytics sx={{ fontSize: 40, color: theme.palette.primary.main }} />,
      title: 'Analytics & Reports',
      description: 'Get detailed insights into your inventory performance and usage patterns.',
    },
    {
      icon: <Security sx={{ fontSize: 40, color: theme.palette.primary.main }} />,
      title: 'Secure & Reliable',
      description: 'Enterprise-grade security with role-based access control and data protection.',
    },
    {
      icon: <Speed sx={{ fontSize: 40, color: theme.palette.primary.main }} />,
      title: 'Fast & Responsive',
      description: 'Lightning-fast performance with a modern, intuitive user interface.',
    },
    {
      icon: <CloudSync sx={{ fontSize: 40, color: theme.palette.primary.main }} />,
      title: 'Cloud-Based',
      description: 'Access your inventory data from anywhere, anytime with cloud synchronization.',
    },
    {
      icon: <DashboardIcon sx={{ fontSize: 40, color: theme.palette.primary.main }} />,
      title: 'Smart Dashboard',
      description: 'Comprehensive dashboard with real-time metrics and actionable insights.',
    },
  ];

  const benefits = [
    'Reduce stockouts and overstock situations',
    'Improve inventory turnover rates',
    'Streamline procurement processes',
    'Enhance decision-making with data-driven insights',
    'Save time with automated inventory tracking',
    'Scale your business with confidence',
  ];

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'background.default' }}>
      {/* Navigation Bar */}
      <AppBar position="static" elevation={0} sx={{ bgcolor: 'white', color: 'text.primary' }}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontWeight: 'bold', color: theme.palette.primary.main }}>
            PantryOS
          </Typography>
          <Stack direction="row" spacing={2}>
            <Button
              color="inherit"
              startIcon={<Login />}
              onClick={() => navigate('/auth/login')}
              sx={{ textTransform: 'none' }}
            >
              Login
            </Button>
            <Button
              variant="contained"
              startIcon={<PersonAdd />}
              onClick={() => navigate('/auth/register')}
              sx={{ textTransform: 'none' }}
            >
              Get Started
            </Button>
          </Stack>
        </Toolbar>
      </AppBar>

      {/* Hero Section */}
      <Container maxWidth="lg" sx={{ py: 8 }}>
        <Box sx={{ 
          display: 'grid', 
          gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, 
          gap: 6, 
          alignItems: 'center' 
        }}>
          <Box>
            <Typography
              variant="h2"
              component="h1"
              gutterBottom
              sx={{
                fontWeight: 'bold',
                background: `linear-gradient(45deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                backgroundClip: 'text',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                fontSize: isMobile ? '2.5rem' : '3.5rem',
              }}
            >
              Smart Inventory Management
            </Typography>
            <Typography variant="h5" color="text.secondary" paragraph sx={{ mb: 4 }}>
              Streamline your inventory operations with our powerful, cloud-based management system. 
              Track, analyze, and optimize your stock levels with ease.
            </Typography>
            <Stack direction="row" spacing={2} sx={{ mb: 4 }}>
              <Button
                variant="contained"
                size="large"
                startIcon={<PersonAdd />}
                onClick={() => navigate('/auth/register')}
                sx={{ 
                  textTransform: 'none', 
                  px: 4, 
                  py: 1.5,
                  fontSize: '1.1rem',
                  borderRadius: 2,
                }}
              >
                Start Free Trial
              </Button>
              <Button
                variant="outlined"
                size="large"
                startIcon={<Login />}
                onClick={() => navigate('/auth/login')}
                sx={{ 
                  textTransform: 'none', 
                  px: 4, 
                  py: 1.5,
                  fontSize: '1.1rem',
                  borderRadius: 2,
                }}
              >
                Sign In
              </Button>
            </Stack>
            <Chip
              icon={<TrendingUp />}
              label="Trusted by 1000+ businesses"
              variant="outlined"
              sx={{ borderRadius: 2 }}
            />
          </Box>
          <Box>
            <Paper
              elevation={8}
              sx={{
                p: 3,
                borderRadius: 3,
                background: `linear-gradient(135deg, ${theme.palette.primary.light}15, ${theme.palette.secondary.light}15)`,
                border: `1px solid ${theme.palette.primary.light}30`,
              }}
            >
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                  Dashboard Overview
                </Typography>
                <Box sx={{ display: 'flex', gap: 1 }}>
                  <Notifications sx={{ color: 'success.main' }} />
                  <Assessment sx={{ color: 'info.main' }} />
                  <Settings sx={{ color: 'warning.main' }} />
                </Box>
              </Box>
              <Box sx={{ 
                display: 'grid', 
                gridTemplateColumns: '1fr 1fr', 
                gap: 2 
              }}>
                <Card sx={{ bgcolor: 'success.light', color: 'white' }}>
                  <CardContent sx={{ p: 2, textAlign: 'center' }}>
                    <Typography variant="h4" sx={{ fontWeight: 'bold' }}>1,247</Typography>
                    <Typography variant="body2">Items in Stock</Typography>
                  </CardContent>
                </Card>
                <Card sx={{ bgcolor: 'warning.light', color: 'white' }}>
                  <CardContent sx={{ p: 2, textAlign: 'center' }}>
                    <Typography variant="h4" sx={{ fontWeight: 'bold' }}>23</Typography>
                    <Typography variant="body2">Low Stock Items</Typography>
                  </CardContent>
                </Card>
                <Card sx={{ bgcolor: 'info.light', color: 'white' }}>
                  <CardContent sx={{ p: 2, textAlign: 'center' }}>
                    <Typography variant="h4" sx={{ fontWeight: 'bold' }}>89%</Typography>
                    <Typography variant="body2">Stock Accuracy</Typography>
                  </CardContent>
                </Card>
                <Card sx={{ bgcolor: 'primary.light', color: 'white' }}>
                  <CardContent sx={{ p: 2, textAlign: 'center' }}>
                    <Typography variant="h4" sx={{ fontWeight: 'bold' }}>156</Typography>
                    <Typography variant="body2">Orders Today</Typography>
                  </CardContent>
                </Card>
              </Box>
            </Paper>
          </Box>
        </Box>
      </Container>

      {/* Features Section */}
      <Box sx={{ bgcolor: 'grey.50', py: 8 }}>
        <Container maxWidth="lg">
          <Typography
            variant="h3"
            component="h2"
            align="center"
            gutterBottom
            sx={{ fontWeight: 'bold', mb: 6 }}
          >
            Powerful Features
          </Typography>
          <Box sx={{ 
            display: 'grid', 
            gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' }, 
            gap: 4 
          }}>
            {features.map((feature, index) => (
              <Card
                key={index}
                sx={{
                  height: '100%',
                  p: 3,
                  textAlign: 'center',
                  transition: 'transform 0.2s, box-shadow 0.2s',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: theme.shadows[8],
                  },
                }}
              >
                <Box sx={{ mb: 2 }}>{feature.icon}</Box>
                <Typography variant="h6" component="h3" gutterBottom sx={{ fontWeight: 'bold' }}>
                  {feature.title}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {feature.description}
                </Typography>
              </Card>
            ))}
          </Box>
        </Container>
      </Box>

      {/* Benefits Section */}
      <Container maxWidth="lg" sx={{ py: 8 }}>
        <Box sx={{ 
          display: 'grid', 
          gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, 
          gap: 6, 
          alignItems: 'center' 
        }}>
          <Box>
            <Typography variant="h3" component="h2" gutterBottom sx={{ fontWeight: 'bold' }}>
              Why Choose PantryOS?
            </Typography>
            <Typography variant="body1" color="text.secondary" paragraph sx={{ mb: 4 }}>
              Our inventory management system is designed to help businesses of all sizes 
              optimize their operations and drive growth through intelligent stock management.
            </Typography>
            <List>
              {benefits.map((benefit, index) => (
                <ListItem key={index} sx={{ px: 0 }}>
                  <ListItemIcon>
                    <Box
                      sx={{
                        width: 8,
                        height: 8,
                        borderRadius: '50%',
                        bgcolor: theme.palette.primary.main,
                      }}
                    />
                  </ListItemIcon>
                  <ListItemText primary={benefit} />
                </ListItem>
              ))}
            </List>
          </Box>
          <Box>
            <Paper
              elevation={4}
              sx={{
                p: 4,
                borderRadius: 3,
                background: `linear-gradient(135deg, ${theme.palette.primary.main}10, ${theme.palette.secondary.main}10)`,
              }}
            >
              <Typography variant="h5" gutterBottom sx={{ fontWeight: 'bold', mb: 3 }}>
                Ready to Get Started?
              </Typography>
              <Typography variant="body1" color="text.secondary" paragraph sx={{ mb: 3 }}>
                Join thousands of businesses that trust PantryOS for their inventory management needs. 
                Start your free trial today and experience the difference.
              </Typography>
              <Stack direction="row" spacing={2}>
                <Button
                  variant="contained"
                  size="large"
                  startIcon={<PersonAdd />}
                  onClick={() => navigate('/auth/register')}
                  sx={{ 
                    textTransform: 'none', 
                    px: 4, 
                    py: 1.5,
                    borderRadius: 2,
                  }}
                >
                  Start Free Trial
                </Button>
                <Button
                  variant="outlined"
                  size="large"
                  startIcon={<Login />}
                  onClick={() => navigate('/auth/login')}
                  sx={{ 
                    textTransform: 'none', 
                    px: 4, 
                    py: 1.5,
                    borderRadius: 2,
                  }}
                >
                  Sign In
                </Button>
              </Stack>
            </Paper>
          </Box>
        </Box>
      </Container>

      {/* Footer */}
      <Box sx={{ bgcolor: 'grey.900', color: 'white', py: 4 }}>
        <Container maxWidth="lg">
          <Box sx={{ 
            display: 'grid', 
            gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, 
            gap: 4 
          }}>
            <Box>
              <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 2 }}>
                PantryOS
              </Typography>
              <Typography variant="body2" color="grey.400">
                Smart inventory management for modern businesses. 
                Streamline your operations and drive growth with our powerful platform.
              </Typography>
            </Box>
            <Box>
              <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 2 }}>
                Quick Links
              </Typography>
              <Stack direction="row" spacing={3}>
                <Button
                  color="inherit"
                  onClick={() => navigate('/auth/login')}
                  sx={{ textTransform: 'none' }}
                >
                  Login
                </Button>
                <Button
                  color="inherit"
                  onClick={() => navigate('/auth/register')}
                  sx={{ textTransform: 'none' }}
                >
                  Register
                </Button>
              </Stack>
            </Box>
          </Box>
          <Divider sx={{ my: 3, borderColor: 'grey.700' }} />
          <Typography variant="body2" color="grey.400" align="center">
            © 2024 PantryOS. All rights reserved.
          </Typography>
        </Container>
      </Box>
    </Box>
  );
};

export default LandingPage; 