import { Metadata } from "next";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { CheckCircle } from "lucide-react";

export const metadata: Metadata = {
  title: "Verification Successful",
  description: "Your email has been verified successfully",
};

export default function VerificationSuccessPage() {
  return (
    <div className="container flex h-screen w-screen flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <div className="flex flex-col items-center space-y-4 text-center">
          <CheckCircle className="h-16 w-16 text-green-500" />
          <h1 className="text-2xl font-semibold tracking-tight">
            Email Verified Successfully
          </h1>
          <p className="text-sm text-muted-foreground">
            Your email has been verified. You can now log in to your account.
          </p>
        </div>
        <div className="flex flex-col space-y-4">
          <Link href="/auth/login">
            <Button className="w-full">Go to Login</Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
