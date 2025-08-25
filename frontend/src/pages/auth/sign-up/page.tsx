import { useFormik } from "formik";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import * as yup from "yup";

import {
  Alert,
  AlertTitle,
  Box,
  Button,
  capitalize,
  Divider,
  FormControl,
  FormLabel,
  Input,
  Paper,
  Tooltip,
  Typography,
} from "@mui/material";

import Logo from "@/components/logo/logo";
import { DEFAULTS } from "@/config";
import NiCheck from "@/icons/nexture/ni-check";
import NiCross from "@/icons/nexture/ni-cross";
import NiCrossSquare from "@/icons/nexture/ni-cross-square";
import { cn } from "@/lib/utils";

const validationSchema = yup.object({
  name: yup.string().required("The field is required").min(3, "Should be at least 3 characters"),
  email: yup.string().required("The field is required").email("Enter a valid email"),
  company: yup.string().required("The field is required").min(3, "Should be at least 3 characters"),
  password: yup
    .string()
    .required("The field is required")
    .min(8, "Should be at least 8 characters")
    .test("uppercase", "Should be an uppercase and a lowercase letter", (value) => {
      const hasUpperCase = /[A-Z]/.test(value);
      const hasLowerCase = /[a-z]/.test(value);
      return hasUpperCase && hasLowerCase;
    })
    .test("symbol", "Should be a special character", (value) => {
      const hasSymbol = /[^A-Za-z 0-9]/g.test(value);
      return hasSymbol;
    }),
});

type InputErrorProps = {
  title: string;
};

const InputErrorTooltip = ({ title }: InputErrorProps) => {
  return (
    <Box className="relative">
      <Tooltip title={title} arrow className="absolute -top-1.5">
        <Button
          startIcon={<NiCrossSquare size="small" />}
          color="error"
          size="small"
          className="group icon-only bg-transparent! outline-0!"
        ></Button>
      </Tooltip>
    </Box>
  );
};

export default function Page() {
  const navigate = useNavigate();
  const [submitted, setSubmitted] = useState(false);

  const formik = useFormik({
    initialValues: {
      name: "",
      email: "",
      company: "",
      password: "",
    },
    validationSchema,
    onSubmit: (values) => {
      console.log(JSON.stringify(values, null, 2));
      navigate(DEFAULTS.appRoot);
    },
    validateOnBlur: false,
    validateOnMount: false,
  });

  const isPasswordLengthValid = () => {
    return formik.values.password.length >= 8;
  };

  const isPasswordCaseValid = () => {
    const hasUpperCase = /[A-Z]/.test(formik.values.password);
    const hasLowerCase = /[a-z]/.test(formik.values.password);
    return hasUpperCase && hasLowerCase;
  };

  const isPasswordSymbolValid = () => {
    const hasSymbol = /[^A-Za-z 0-9]/g.test(formik.values.password);
    return hasSymbol;
  };

  const googleSVG = () => {
    return (
      <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M19.6169 10.2876C19.6169 9.60932 19.5561 8.95714 19.443 8.33105H10.4343V12.0354H15.5822C15.3561 13.2267 14.6778 14.2354 13.6604 14.9137V17.3224H16.7648C18.5735 15.6528 19.6169 13.2006 19.6169 10.2876Z"
          fill="#4285F4"
        />
        <path
          d="M10.4346 19.6346C13.0172 19.6346 15.1825 18.7825 16.7651 17.3216L13.6607 14.9129C12.8086 15.4868 11.7216 15.8346 10.4346 15.8346C7.94768 15.8346 5.83464 14.1564 5.07812 11.8955H1.89551V14.3651C3.46942 17.4868 6.69551 19.6346 10.4346 19.6346Z"
          fill="#34A853"
        />
        <path
          d="M5.07832 11.8866C4.88702 11.3127 4.77398 10.704 4.77398 10.0692C4.77398 9.4344 4.88702 8.8257 5.07832 8.25179V5.78223H1.89572C1.24354 7.06918 0.869629 8.52136 0.869629 10.0692C0.869629 11.617 1.24354 13.0692 1.89572 14.3561L4.37398 12.4257L5.07832 11.8866Z"
          fill="#FBBC05"
        />
        <path
          d="M10.4346 4.31358C11.8433 4.31358 13.0955 4.80054 14.0955 5.73967L16.8346 3.00054C15.1738 1.45271 13.0172 0.504883 10.4346 0.504883C6.69551 0.504883 3.46942 2.65271 1.89551 5.78314L5.07812 8.25271C5.83464 5.99184 7.94768 4.31358 10.4346 4.31358Z"
          fill="#EA4335"
        />
      </svg>
    );
  };

  const githubSVG = () => {
    return (
      <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M10 0.0693359C4.475 0.0693359 0 4.54434 0 10.0693C0 14.4943 2.8625 18.2318 6.8375 19.5568C7.3375 19.6443 7.525 19.3443 7.525 19.0818C7.525 18.8443 7.5125 18.0568 7.5125 17.2193C5 17.6818 4.35 16.6068 4.15 16.0443C4.0375 15.7568 3.55 14.8693 3.125 14.6318C2.775 14.4443 2.275 13.9818 3.1125 13.9693C3.9 13.9568 4.4625 14.6943 4.65 14.9943C5.55 16.5068 6.9875 16.0818 7.5625 15.8193C7.65 15.1693 7.9125 14.7318 8.2 14.4818C5.975 14.2318 3.65 13.3693 3.65 9.54434C3.65 8.45684 4.0375 7.55684 4.675 6.85684C4.575 6.60684 4.225 5.58184 4.775 4.20684C4.775 4.20684 5.6125 3.94434 7.525 5.23184C8.325 5.00684 9.175 4.89434 10.025 4.89434C10.875 4.89434 11.725 5.00684 12.525 5.23184C14.4375 3.93184 15.275 4.20684 15.275 4.20684C15.825 5.58184 15.475 6.60684 15.375 6.85684C16.0125 7.55684 16.4 8.44434 16.4 9.54434C16.4 13.3818 14.0625 14.2318 11.8375 14.4818C12.2 14.7943 12.5125 15.3943 12.5125 16.3318C12.5125 17.6693 12.5 18.7443 12.5 19.0818C12.5 19.3443 12.6875 19.6568 13.1875 19.5568C17.1375 18.2318 20 14.4818 20 10.0693C20 4.54434 15.525 0.0693359 10 0.0693359Z"
          fill="#1B1F23"
        />
      </svg>
    );
  };

  return (
    <Box className="flex min-h-screen w-full items-center justify-center p-4">
      <Paper elevation={3} className="bg-background-paper shadow-darker-xs w-[32rem] max-w-full rounded-4xl py-14">
        <Box className="flex flex-col gap-4 px-8 sm:px-14">
          <Box className="flex flex-col">
            <Box className="mb-14 flex justify-center">
              <Logo classNameMobile="hidden" />
            </Box>

            <Box className="flex flex-col gap-10">
              <Box className="flex flex-col">
                <Typography variant="h1" component="h1" className="mb-2">
                  Sign up
                </Typography>
                <Typography variant="body1" className="text-text-primary">
                  Create your account in just a few steps and start today.
                </Typography>
              </Box>

              <Box className="flex flex-col gap-5">
                <Box className="flex flex-col gap-2 md:flex-row">
                  <Button variant="outlined" color="grey" className="flex-none md:w-1/2">
                    <Box className="mr-2">{googleSVG()}</Box>Sign in with Google
                  </Button>
                  <Button variant="outlined" color="grey" className="flex-none md:w-1/2">
                    <Box className="mr-2">{githubSVG()}</Box>Sign in with Google
                  </Button>
                </Box>

                <Divider className="text-text-secondary my-0 text-sm">OR</Divider>

                <Box
                  component={"form"}
                  onSubmit={(event) => {
                    setSubmitted(true);
                    formik.handleSubmit(event);
                  }}
                  className="flex flex-col"
                >
                  <FormControl className="outlined" variant="standard" size="small">
                    <FormLabel component="label" className="flex flex-row">
                      Name{" "}
                      {formik.touched.name && formik.errors.name && <InputErrorTooltip title={formik.errors.name} />}
                    </FormLabel>
                    <Input
                      id="name"
                      name="name"
                      placeholder=""
                      value={formik.values.name}
                      onChange={formik.handleChange}
                      onBlur={formik.handleBlur}
                    />
                  </FormControl>

                  <FormControl className="outlined" variant="standard" size="small">
                    <FormLabel component="label" className="flex flex-row">
                      Email{" "}
                      {formik.touched.email && formik.errors.email && <InputErrorTooltip title={formik.errors.email} />}
                    </FormLabel>
                    <Input
                      id="email"
                      name="email"
                      placeholder=""
                      value={formik.values.email}
                      onChange={formik.handleChange}
                      onBlur={formik.handleBlur}
                    />
                  </FormControl>

                  <FormControl className="outlined" variant="standard" size="small">
                    <FormLabel component="label" className="flex flex-row">
                      Company{" "}
                      {formik.touched.company && formik.errors.company && (
                        <InputErrorTooltip title={formik.errors.company} />
                      )}
                    </FormLabel>
                    <Input
                      id="company"
                      name="company"
                      placeholder=""
                      value={formik.values.company}
                      onChange={formik.handleChange}
                      onBlur={formik.handleBlur}
                    />
                  </FormControl>

                  <FormControl className="outlined" variant="standard" size="small">
                    <FormLabel component="label" className="flex flex-row">
                      Password{" "}
                      {formik.touched.password && formik.errors.password && (
                        <InputErrorTooltip title={formik.errors.password} />
                      )}
                    </FormLabel>
                    <Input
                      id="password"
                      name="password"
                      placeholder=""
                      autoComplete="off"
                      value={formik.values.password}
                      onChange={formik.handleChange}
                      onBlur={formik.handleBlur}
                    />
                    <Typography variant="body2" className="text-text-secondary mt-2 inline-block align-middle">
                      <span className="inline">Must be</span>
                      <span
                        className={cn(
                          "mx-1 inline-block h-4 w-4 rounded-md align-text-bottom",
                          isPasswordLengthValid() ? "bg-success text-text-contrast" : "bg-grey-100 text-text-secondary",
                        )}
                      >
                        {isPasswordLengthValid() ? (
                          <NiCheck size={"tiny"}></NiCheck>
                        ) : (
                          <NiCross size={"tiny"}></NiCross>
                        )}
                      </span>
                      <span className={cn("inline font-semibold", isPasswordLengthValid() && "text-success")}>
                        at least 8 characters long,{" "}
                      </span>
                      <span className="inline">must contain</span>
                      <span
                        className={cn(
                          "mx-1 inline-block h-4 w-4 rounded-md align-text-bottom",
                          isPasswordCaseValid() ? "bg-success text-text-contrast" : "bg-grey-100 text-text-secondary",
                        )}
                      >
                        {isPasswordCaseValid() ? <NiCheck size={"tiny"}></NiCheck> : <NiCross size={"tiny"}></NiCross>}
                      </span>
                      <span className={cn("inline font-semibold", isPasswordCaseValid() && "text-success")}>
                        lowercase and uppercase letters,{" "}
                      </span>
                      <span className="inline">must have at least</span>
                      <span
                        className={cn(
                          "mx-1 inline-block h-4 w-4 rounded-md align-text-bottom",
                          isPasswordSymbolValid() ? "bg-success text-text-contrast" : "bg-grey-100 text-text-secondary",
                        )}
                      >
                        {isPasswordSymbolValid() ? (
                          <NiCheck size={"tiny"}></NiCheck>
                        ) : (
                          <NiCross size={"tiny"}></NiCross>
                        )}
                      </span>
                      <span className={cn("inline font-semibold", isPasswordSymbolValid() && "text-success")}>
                        one special character.
                      </span>
                    </Typography>
                  </FormControl>
                  {submitted && !formik.isValid && (
                    <Alert severity="error" icon={<NiCrossSquare />} className="neutral bg-background-paper/60! mb-4">
                      <AlertTitle variant="subtitle2">The following inputs have errors!</AlertTitle>
                      {Object.entries(formik.errors).map(([key, value]) => {
                        return (
                          <Box className="flex flex-row gap-0.5" key={crypto.randomUUID()}>
                            <Typography variant="body2" className="text-error">
                              {capitalize(key)}:
                            </Typography>
                            <Typography variant="body2" className="text-text-primary">
                              {value}
                            </Typography>
                          </Box>
                        );
                      })}
                    </Alert>
                  )}
                  <Box className="flex flex-col gap-2">
                    <Link
                      to="/auth/password-reset"
                      className="link-text-secondary link-underline-hover text-center text-sm font-semibold"
                    >
                      Reset Password
                    </Link>
                    <Button type="submit" variant="contained" className="mb-4">
                      Continue
                    </Button>
                  </Box>

                  <Typography variant="body2" className="text-text-secondary">
                    By clicking Continue, Sign in with Google, or Sign in with GitHub, you agree to the{" "}
                    <Link target="_blank" to="/auth/terms-and-conditions" className="link-primary link-underline-hover">
                      Terms and Conditions
                    </Link>{" "}
                    and{" "}
                    <Link target="_blank" to="/auth/privacy-policy" className="link-primary link-underline-hover">
                      Privacy Policy
                    </Link>
                    .
                  </Typography>
                </Box>
              </Box>
              <Divider className="text-text-secondary my-0 text-sm"></Divider>
              <Box className="flex flex-col">
                <Typography variant="h5" component="h5">
                  Sign in
                </Typography>
                <Typography variant="body1" className="text-text-secondary">
                  If you already have an account, please{" "}
                  <Link to="/auth/sign-in" className="link-primary link-underline-hover">
                    sign in
                  </Link>
                  .
                </Typography>
              </Box>
            </Box>
          </Box>
        </Box>
      </Paper>
    </Box>
  );
}
