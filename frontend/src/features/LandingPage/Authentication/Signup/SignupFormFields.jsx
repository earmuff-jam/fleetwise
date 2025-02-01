import { produce } from 'immer';
import { useDispatch } from 'react-redux';
import { Stack, TextField, InputAdornment, Typography } from '@mui/material';
import { authActions } from '@features/LandingPage/authSlice';
import SignupHelperText from '@features/LandingPage/Authentication/Signup/SignupHelperText';

export default function SignupFormFields({ formFields, setFormFields, isValidUserEmail, loading, submit }) {
  const dispatch = useDispatch();

  const handleInput = (event) => {
    const { name, value } = event.target;
    setFormFields(
      produce(formFields, (draft) => {
        draft[name].value = value;
        draft[name].errorMsg = '';

        for (const validator of draft[name].validators) {
          if (validator.validate(value)) {
            draft[name].errorMsg = validator.message;
            break;
          }
        }
      })
    );
  };

  return (
    <Stack spacing={1}>
      <Typography variant="subtitle2" color="text.secondary">
        {formFields['username'].label} {formFields['username'].required && '*'}
      </Typography>
      <TextField
        id={formFields['username'].name}
        name={formFields['username'].name}
        value={formFields['username'].value}
        size={formFields['username'].size}
        type={formFields['username'].type}
        variant={formFields['username'].variant}
        autoComplete={formFields['username'].autocomplete}
        placeholder={formFields['username'].placeholder}
        onChange={handleInput}
        required={formFields['username'].required}
        fullWidth={formFields['username'].fullWidth}
        error={!!formFields['username'].errorMsg}
        helperText={formFields['username'].errorMsg}
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['username'].icon}</InputAdornment>,
        }}
      />

      <Typography variant="subtitle2" color="text.secondary">
        {formFields['email'].label} {formFields['email'].required && '*'}
      </Typography>
      <TextField
        id={formFields['email'].name}
        name={formFields['email'].name}
        value={formFields['email'].value}
        size={formFields['email'].size}
        type={formFields['email'].type}
        variant={formFields['email'].variant}
        autoComplete={formFields['email'].autocomplete}
        placeholder={formFields['email'].placeholder}
        onChange={handleInput}
        required={formFields['email'].required}
        fullWidth={formFields['email'].fullWidth}
        error={!!formFields['email'].errorMsg}
        helperText={
          formFields['email'].errorMsg || <SignupHelperText isEmailUnique={isValidUserEmail} loading={loading} />
        }
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        onBlur={() =>
          formFields.email.value.length > 3 && dispatch(authActions.isValidUserEmail({ email: formFields.email.value }))
        }
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['email'].icon}</InputAdornment>,
        }}
      />

      <Typography variant="subtitle2" color="text.secondary">
        {formFields['password'].label} {formFields['password'].required && '*'}
      </Typography>
      <TextField
        id={formFields['password'].name}
        name={formFields['password'].name}
        value={formFields['password'].value}
        size={formFields['password'].size}
        type={formFields['password'].type}
        variant={formFields['password'].variant}
        autoComplete={formFields['password'].autocomplete}
        placeholder={formFields['password'].placeholder}
        onChange={handleInput}
        required={formFields['password'].required}
        fullWidth={formFields['password'].fullWidth}
        error={!!formFields['password'].errorMsg}
        helperText={formFields['password'].errorMsg}
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['password'].icon}</InputAdornment>,
        }}
      />

      <Typography variant="subtitle2" color="text.secondary">
        {formFields['birthday'].label} {formFields['birthday'].required && '*'}
      </Typography>
      <TextField
        id={formFields['birthday'].name}
        name={formFields['birthday'].name}
        value={formFields['birthday'].value}
        type={formFields['birthday'].type}
        size={formFields['birthday'].size}
        variant={formFields['birthday'].variant}
        autoComplete={formFields['birthday'].autocomplete}
        placeholder={formFields['birthday'].placeholder}
        onChange={handleInput}
        required={formFields['birthday'].required}
        fullWidth={formFields['birthday'].fullWidth}
        error={!!formFields['birthday'].errorMsg}
        helperText={formFields['birthday'].errorMsg}
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['birthday'].icon}</InputAdornment>,
        }}
      />
    </Stack>
  );
}
