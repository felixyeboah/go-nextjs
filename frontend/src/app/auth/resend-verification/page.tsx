import { Metadata } from "next";
import Link from "next/link";
import { ResendVerificationForm } from "@/components/auth/resend-verification-form";

export const metadata: Metadata = {
  title: "Resend Verification",
  description: "Resend verification email to your account",
};

export default function ResendVerificationPage() {
  return (
    <div className="container flex h-screen w-screen flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">
            Resend Verification Email
          </h1>
          <p className="text-sm text-muted-foreground">
            Enter your email to receive a new verification link
          </p>
        </div>
        <ResendVerificationForm />
        <p className="px-8 text-center text-sm text-muted-foreground">
          <Link
            href="/auth/login"
            className="hover:text-brand underline underline-offset-4"
          >
            Back to login
          </Link>
        </p>
      </div>
    </div>
  );
}
