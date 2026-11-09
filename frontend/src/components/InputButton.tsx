import "./InputButton.css";

type InputButtonProps = {
  text: string;
  className: string;
  onClick: () => void;
};

function InputButton({
  text: ButtonTextString,
  className: className,
  onClick,
}: InputButtonProps) {
  return (
    <button className={className} onClick={onClick}>
      {ButtonTextString}
    </button>
  );
}

export default InputButton;
