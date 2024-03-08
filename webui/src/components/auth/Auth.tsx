// Base Auth component

import { getOAuthUrl } from "@/api/auth";

// Reusable component for triggering oauth flows
type Props = {
  name: string;
};

const handleAuthProvider = async (provider: string) => {
  const oauthUrl = await getOAuthUrl(provider);
  const authUrl = oauthUrl.url;
  window.location.href = authUrl;
};

const Auth = (props: Props) => {
  return (
    <div>
      <button
        className={`auth ${props.name.toLowerCase()}`}
        type="button"
        onClick={() => handleAuthProvider(props.name)}
      >
        Login with {props.name}
      </button>
    </div>
  );
};

export default Auth;
