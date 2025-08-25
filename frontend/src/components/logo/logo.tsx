import { cn } from "@/lib/utils";

export default function Logo() {
  // 1. Define your styles as a JavaScript object
  const logoTextStyle = {
    fontFamily: 'system-ui, -apple-system, "Segoe UI", Roboto, Helvetica, Arial, sans-serif',
    fontSize: '16px',
    fontWeight: 600,
    fill: '#333',
  };

  return (
    <>
      <svg width="200" height="60" viewBox="0 0 200 60" xmlns="http://www.w3.org/2000/svg">
        {/* 2. The <style> tag has been removed */}

        <g transform="translate(10, 8)">
          <path d="M 0,0 L 0,44" stroke="#333" strokeWidth="8" strokeLinecap="round" />
          <path d="M 0,0 A 22,22 0 1 1 0,44" fill="none" stroke="#333" strokeWidth="8" />
          <path d="M 20,22 Q 32,10 45,18" fill="#4CAF50" />
        </g>

        {/* 3. Apply the style object directly to the element */}
        <text x="75" y="42" style={logoTextStyle}>
          Pantry OS
        </text>
      </svg>
    </>
  );
}