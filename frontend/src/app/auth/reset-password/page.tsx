import { Metadata } from "next";
import Link from "next/link";
import { Suspense } from "react";
import ResetPasswordForm from "@/components/auth/reset-password-form";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export const metadata: Metadata = {
  title: "Reset Password",
  description: "Reset your password",
};

export default function ResetPasswordPage() {
  return (
    <div className="container flex h-screen w-screen flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[400px]">
        <Card>
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl font-semibold tracking-tight text-center">
              Reset Password
            </CardTitle>
            <CardDescription className="text-center">
              Enter your new password below
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Suspense fallback={<div>Loading...</div>}>
              <ResetPasswordForm />
            </Suspense>
          </CardContent>
          <CardFooter>
            <div className="text-center text-sm text-muted-foreground w-full">
              <Link
                href="/auth/login"
                className="hover:text-primary underline underline-offset-4"
              >
                Back to login
              </Link>
            </div>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
