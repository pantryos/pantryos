import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Auth() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate("/auth/sign-in");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return <div>Welcome!</div>;
}
