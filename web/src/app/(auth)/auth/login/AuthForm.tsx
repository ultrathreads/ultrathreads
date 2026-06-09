'use client';

import { useState, FormEvent, useCallback } from 'react';
import { toast } from 'sonner';
import { ApiBusinessError } from '@/lib/api/client';

export interface AuthFieldConfig {
  name: string;
  label: string;
  type?: 'text' | 'email' | 'password';
  placeholder?: string;
  autoComplete?: string;
  required?: boolean;
  validate?: (value: string) => string | undefined;
}

interface AuthFormProps {
  title: string;
  subtitle: string;
  fields: AuthFieldConfig[];
  submitLabel: string;
  submittingLabel: string;
  onSubmit: (values: Record<string, string>) => Promise<void>;
  footer?: React.ReactNode;
  renderExtraContent?: (fieldName: string) => React.ReactNode;
}

export default function AuthForm({
  title, subtitle, fields, submitLabel, submittingLabel,
  onSubmit, footer, renderExtraContent,
}: AuthFormProps) {
  const [values, setValues] = useState<Record<string, string>>(
    () => Object.fromEntries(fields.map((f) => [f.name, '']))
  );
  const [fieldErrors, setFieldErrors] = useState<Record<string, string | undefined>>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  const validate = useCallback((): Record<string, string> => {
    const errors: Record<string, string> = {};
    for (const field of fields) {
      const val = values[field.name]?.trim() ?? '';
      if (field.required && !val) {
        errors[field.name] = `请输入${field.label}`;
      } else if (field.validate) {
        const err = field.validate(val);
        if (err) errors[field.name] = err;
      }
    }
    return errors;
  }, [fields, values]);

  const canSubmit = !isSubmitting && fields.every((f) => {
    const val = values[f.name]?.trim() ?? '';
    if (f.required && !val) return false;
    if (f.validate && f.validate(val)) return false;
    return true;
  });

  const handleChange = (name: string, val: string) => {
    setValues((prev) => ({ ...prev, [name]: val }));
    setServerError(null);
    if (fieldErrors[name]) {
      setFieldErrors((prev) => ({ ...prev, [name]: undefined }));
    }
  };

  const handleBlur = (name: string) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
    const field = fields.find((f) => f.name === name);
    if (!field) return;
    const val = values[name]?.trim() ?? '';
    let err: string | undefined;
    if (field.required && !val) {
      err = `请输入${field.label}`;
    } else if (field.validate) {
      err = field.validate(val);
    }
    setFieldErrors((prev) => ({ ...prev, [name]: err }));
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setServerError(null);
    setTouched(Object.fromEntries(fields.map((f) => [f.name, true])));

    const errors = validate();
    setFieldErrors(errors);
    if (Object.keys(errors).length > 0) return;

    setIsSubmitting(true);
    try {
      await onSubmit(values);
    } catch (err) {
      if (err instanceof ApiBusinessError) {
        setServerError(err.message);
      } else {
        const msg = err instanceof Error ? err.message : '未知错误';
        toast.error(msg);
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-header">
        <h1 className="auth-title">{title}</h1>
        <p className="auth-subtitle">{subtitle}</p>
      </div>

      {serverError && (
        <div className="auth-error" role="alert">{serverError}</div>
      )}

      <form className="auth-form" onSubmit={handleSubmit} noValidate>
        {fields.map((field) => {
          const error = fieldErrors[field.name];
          const showError = touched[field.name] && !!error;
          const errorId = `${field.name}-error`;

          return (
            <div className="form-group" key={field.name}>
              <label className="form-label" htmlFor={field.name}>
                {field.label}
              </label>
              <input
                id={field.name}
                type={field.type ?? 'text'}
                className={`form-input ${showError ? 'form-error' : ''}`}
                placeholder={field.placeholder}
                autoComplete={field.autoComplete}
                value={values[field.name]}
                onChange={(e) => handleChange(field.name, e.target.value)}
                onBlur={() => handleBlur(field.name)}
                aria-invalid={showError}
                aria-describedby={showError ? errorId : undefined}
              />
              {showError && (
                <p id={errorId} className="form-hint form-hint--error" role="alert">
                  {error}
                </p>
              )}
              {renderExtraContent?.(field.name)}
            </div>
          );
        })}

        <div className="auth-actions">
          <button
            type="submit"
            className="auth-btn"
            disabled={!canSubmit}
            aria-disabled={!canSubmit}
          >
            {isSubmitting ? submittingLabel : submitLabel}
          </button>
        </div>
      </form>

      {footer && <div className="auth-footer">{footer}</div>}
    </div>
  );
}