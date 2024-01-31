function updateCheckboxes(selectElement) {
  // Get the selected value from the dropdown
  var SV = selectElement.value[0];
  var parent = selectElement.parentNode;
  parent.querySelector('label[for="type1"]').innerText = SV;
  parent.querySelector('label[for="type2"]').innerText = "T" + SV;
  parent.querySelector('label[for="type3"]').innerText = "T" + SV + SV;
}
